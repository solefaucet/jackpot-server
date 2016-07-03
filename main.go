package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
	s "github.com/solefaucet/jackpot-server/services/storage"
	"github.com/solefaucet/jackpot-server/services/storage/mysql"
	w "github.com/solefaucet/jackpot-server/services/wallet"
	"github.com/solefaucet/jackpot-server/services/wallet/core"
	"github.com/solefaucet/jackpot-server/utils"
	grayloghook "github.com/yumimobi/logrus-graylog2-hook"
)

var (
	logger  = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)
	wallet  w.Wallet
	storage s.Storage
)

func initService() {
	// configuration
	initConfig()

	// log
	l := utils.Must(logrus.ParseLevel(config.Log.Level)).(logrus.Level)
	logrus.SetLevel(l)
	logrus.SetOutput(os.Stdout)

	// logging hooks
	graylogHookLevelThreshold := utils.Must(logrus.ParseLevel(config.Log.Graylog.Level)).(logrus.Level)
	graylogHook := utils.Must(
		grayloghook.New(
			config.Log.Graylog.Address,
			config.Log.Graylog.Facility,
			map[string]interface{}{
				"go_version": goVersion,
				"build_time": buildTime,
				"git_commit": gitCommit,
			},
			graylogHookLevelThreshold,
		),
	).(logrus.Hook)
	logrus.AddHook(graylogHook)

	// storage
	store := mysql.New(config.DB.DataSourceName)
	store.SetMaxOpenConns(config.DB.MaxOpenConns)
	store.SetMaxIdleConns(config.DB.MaxIdleConns)
	storage = store

	// wallet
	wallet = utils.Must(
		core.New(
			config.Wallet.Host,
			config.Wallet.Username,
			config.Wallet.Password,
		),
	).(w.Wallet)

	// init cronjob
	initCronjob()
}

func main() {
	initService()

	// on service stop, log and maybe do some cleanup jobs
	onServiceStop := func() {
		logrus.WithFields(logrus.Fields{
			"event": models.LogEventServiceStateChanged,
		}).Info("service is stopping...")
	}
	defer onServiceStop()
	go catch(onServiceStop)

	logrus.WithFields(logrus.Fields{
		"event": models.LogEventServiceStateChanged,
	}).Info("service up")
}

func initCronjob() {
	c := cron.New()
	utils.Must(nil, c.AddFunc("@every 1m", safeFuncWrapper(addNewBlock))) // add new block every 1 minute
	c.Start()
}

func addNewBlock() {
	entry := logrus.WithField("event", models.LogEventAddNewBlock)

	// get latest block from db
	block, err := storage.GetLatestBlock()
	if err != nil && err != jerrors.ErrNotFound {
		entry.WithError(err).Error("fail to get latest block from database")
		return
	}

	// get block from blockchain
	bestBlock := err == jerrors.ErrNotFound
	blockchainBlock, err := wallet.GetBlock(bestBlock, block.Height+1)
	switch {
	case err == jerrors.ErrNoNewBlock:
		entry.Info("no new block ahead")
		return
	case err != nil:
		entry.WithError(err).Error("fail to get block from blockchain")
		return
	}

	// save block to storage
	block = models.Block{
		Hash:           blockchainBlock.Hash,
		Height:         blockchainBlock.Height,
		BlockCreatedAt: blockchainBlock.BlockCreatedAt.UTC(),
	}
	if err := storage.SaveBlock(block); err != nil {
		entry.WithError(err).Error("fail to save block to database")
		return
	}

	entry.Info("add new block successfully")

	// keep adding until no blocks is ahead
	addNewBlock()
}

func catch(then func()) {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
	if then != nil {
		then()
		os.Exit(1)
	}
}

// wrap a function with recover
func safeFuncWrapper(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				logrus.WithFields(logrus.Fields{
					"error": fmt.Sprintf("%v", err),
					"stack": string(buf[:n]),
				}).Error("panic")
				logger.Printf("%v\n%s\n", err, buf)
			}
		}()
		f()
	}
}
