package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"github.com/solefaucet/jackpot-server/models"
	"github.com/solefaucet/jackpot-server/utils"
	grayloghook "github.com/yumimobi/logrus-graylog2-hook"
)

var (
	logger  = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)
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

	entry.Info("add new block successfully")
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
