package main

import (
	"reflect"

	"gopkg.in/go-playground/validator.v8"

	"github.com/go-sql-driver/mysql"
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
	} `validate:"required"`
	Wallet struct {
		Host     string `validate:"required"`
		Username string `validate:"required"`
		Password string `validate:"required"`
	} `validate:"required"`
	DB struct {
		DataSourceName string `validate:"required,dsn"`
		MaxOpenConns   int    `validate:"required,min=1"`
		MaxIdleConns   int    `validate:"required,min=1,ltefield=MaxOpenConns"`
	} `validate:"required"`
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

	config.Wallet.Host = viper.GetString("wallet_rpc_host")
	config.Wallet.Username = viper.GetString("wallet_rpc_username")
	config.Wallet.Password = viper.GetString("wallet_rpc_password")

	config.DB.DataSourceName = viper.GetString("db_dsn")
	config.DB.MaxOpenConns = viper.GetInt("db_max_open_conns")
	config.DB.MaxIdleConns = viper.GetInt("db_max_idle_conns")

	// validate config
	utils.Must(nil, validateConfiguration(config))
}

func validateConfiguration(c configuration) error {
	validate := validator.New(&validator.Config{TagName: "validate"})
	utils.Must(nil, validate.RegisterValidation("dsn", dsnValidator))
	return validate.Struct(c)
}

func dsnValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	dsn, err := mysql.ParseDSN(field.String())
	return err == nil && dsn.ParseTime
}
