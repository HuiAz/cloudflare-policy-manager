package config

import "github.com/spf13/viper"

type CloudflareConfig struct {
	CloudflareAccountId string
	CloudflareApiKey    string
	CloudflareEmail     string
	Devices             []string
	Mode                string
}

type Config struct {
	CloudflareConfig
	Rules []string
}

var Conf Config

func Parse(configPath *string) {
	// read from config yaml file
	viper.SetConfigFile(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file, err: " + err.Error())
	}
	Conf.CloudflareAccountId = viper.GetString("cloudflare.account_id")
	Conf.CloudflareApiKey = viper.GetString("cloudflare.api_key")
	Conf.CloudflareEmail = viper.GetString("cloudflare.email")
	Conf.Mode = viper.GetString("cloudflare.mode")
	Conf.Devices = viper.GetStringSlice("cloudflare.devices")
	Conf.Rules = viper.GetStringSlice("rules")
}
