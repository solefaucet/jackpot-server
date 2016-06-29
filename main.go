package main

import (
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/models"
	"github.com/solefaucet/jackpot-server/utils"
	grayloghook "github.com/yumimobi/logrus-graylog2-hook"
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

func catch(then func()) {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
	if then != nil {
		then()
		os.Exit(1)
	}
}
