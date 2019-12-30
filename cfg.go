package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	URL      string
	Username string
	Password string
}

func parse(config *Config) error {
	pflag.StringVar(&config.URL, "url", "fritz.box", "FritzBox URL")
	pflag.StringVar(&config.Username, "username", "admin", "FritzBox User")
	pflag.StringVar(&config.Password, "password", "admin", "FritzBox Password")
	pflag.Parse()
	return viper.BindPFlags(pflag.CommandLine)
}
