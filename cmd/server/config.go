package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var requiredFields = map[string]string{
	"database.user":     "string",
	"database.password": "string",
	"database.host":     "string",
	"database.name":     "string",
}

func initConfig() error {
	viper.SetDefault("host", "")
	viper.SetDefault("port", "50051")
	viper.SetDefault("tls_cert", "./cert.pem")
	viper.SetDefault("tls_key", "./cert-key.pem")
	viper.SetDefault("access_token_key", "7QC/y4Dkke2izCGyArkfH074ETD9Hyf6PxIV")
	viper.SetDefault("refresh_token_key", "tA2ZFqRCgYBEX4Y9/Q4Au9U0qrbW2oBcqJ8uRPavj9g=")
	viper.SetDefault("database.user", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.name", "")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.nfc-cash-system/server")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("server")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("could not load config: %w", err)
		}
	}

	if err := viper.SafeWriteConfigAs("./config.yml"); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("could not save config: %w", err)
		}
	}

	flag.String("address", "", "Host address for server")
	flag.String("port", viper.GetString("port"), "Host port for server")
	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return fmt.Errorf("could not bind flag values: %v\n", err)
	}

	err = checkRequired()
	if err != nil {
		return err
	}

	return nil
}

func checkRequired() error {
	for name, t := range requiredFields {
		switch t {
		case "string":
			val := viper.GetString(name)
			if val == "" {
				return fmt.Errorf("required setting %s is missing", name)
			}
		}
	}
	return nil
}
