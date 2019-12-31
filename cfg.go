package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	URL string
}

func parse(config *Config) error {
	pflag.StringVar(&config.URL, "url", "fritz.box", "FritzBox URL")
	pflag.Parse()
	return viper.BindPFlags(pflag.CommandLine)
}
