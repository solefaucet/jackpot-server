package main

import (
	"gopkg.in/go-playground/validator.v8"

	"github.com/solefaucet/jackpot-server/utils"
	"github.com/spf13/viper"
)

type graylog struct {
	Address  string `mapstructure:"address" validate:"required"`
	Level    string `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
	Facility string `mapstructure:"facility" validate:"required"`
}

type configuration struct {
	Log struct {
		Level   string  `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
		Graylog graylog `mapstructure:"graylog" validate:"required,dive"`
	} `mapstructure:"log" validate:"required"`
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("jackpot") // will turn into uppercase, e.g. JACKPOT_PORT
	viper.AutomaticEnv()

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default
	config.Log.Level = viper.GetString("log_level")
	config.Log.Graylog.Address = viper.GetString("graylog_address")
	config.Log.Graylog.Level = viper.GetString("graylog_level")
	config.Log.Graylog.Facility = viper.GetString("graylog_facility")

	// validate config
	utils.Must(nil, validateConfiguration(config))
}

func validateConfiguration(c configuration) error {
	validate := validator.New(&validator.Config{TagName: "validate"})
	return validate.Struct(c)
}
