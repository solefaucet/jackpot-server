package main

import (
	"reflect"
	"time"

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
	HTTP struct {
		Address string `validate:"required"`
		Mode    string `validate:"required,eq=release|eq=test|eq=debug"`
	} `validate:"required"`
	Log struct {
		Level   string  `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
		Graylog graylog `mapstructure:"graylog" validate:"required,dive"`
	} `validate:"required"`
	Wallet struct {
		Host            string `validate:"required"`
		Username        string `validate:"required"`
		Password        string `validate:"required"`
		MinConfirms     int64  `validate:"required,min=1"`
		SentFromAccount string `validate:"required"`
	} `validate:"required"`
	Coin struct {
		Type       string `validate:"required"`
		Label      string `validate:"required"`
		TxURL      string `validate:"required"`
		AddressURL string `validate:"required"`
	} `validate:"required"`
	DB struct {
		DataSourceName string `validate:"required,dsn"`
		MaxOpenConns   int    `validate:"required,min=1"`
		MaxIdleConns   int    `validate:"required,min=1,ltefield=MaxOpenConns"`
	} `validate:"required"`
	Jackpot struct {
		DestAddress    string  `validate:"required"`
		TransactionFee float64 `validate:"required,min=0,lt=1"`
		Duration       time.Duration
	} `validate:"required"`
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("jackpot") // will turn into uppercase, e.g. JACKPOT_PORT
	viper.AutomaticEnv()

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default
	config.HTTP.Mode = viper.GetString("mode")
	config.HTTP.Address = viper.GetString("address")

	config.Log.Level = viper.GetString("log_level")
	config.Log.Graylog.Address = viper.GetString("graylog_address")
	config.Log.Graylog.Level = viper.GetString("graylog_level")
	config.Log.Graylog.Facility = viper.GetString("graylog_facility")

	config.Wallet.Host = viper.GetString("wallet_rpc_host")
	config.Wallet.Username = viper.GetString("wallet_rpc_username")
	config.Wallet.Password = viper.GetString("wallet_rpc_password")
	config.Wallet.MinConfirms = int64(viper.GetInt("wallet_min_confirms"))
	config.Wallet.SentFromAccount = viper.GetString("wallet_sent_from_account")

	config.Coin.Type = viper.GetString("coin_type")
	config.Coin.Label = viper.GetString("coin_label")
	config.Coin.TxURL = viper.GetString("coin_tx_url")
	config.Coin.AddressURL = viper.GetString("coin_address_url")

	config.DB.DataSourceName = viper.GetString("db_dsn")
	config.DB.MaxOpenConns = viper.GetInt("db_max_open_conns")
	config.DB.MaxIdleConns = viper.GetInt("db_max_idle_conns")

	config.Jackpot.DestAddress = viper.GetString("dest_address")
	config.Jackpot.TransactionFee = viper.GetFloat64("transaction_fee")
	config.Jackpot.Duration = utils.Must(time.ParseDuration(viper.GetString("duration"))).(time.Duration)

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
