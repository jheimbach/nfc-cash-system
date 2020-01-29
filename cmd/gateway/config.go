package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func initConfig() error {
	viper.SetDefault("rest_port", "8080")
	viper.SetDefault("rest_host", "")
	viper.SetDefault("grpc_port", "50051")
	viper.SetDefault("grpc_host", "")
	viper.SetDefault("tls_cert", "./cert.pem")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.nfc-cash-system/gateway")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("gateway")
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

	return nil
}
