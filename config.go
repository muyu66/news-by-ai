package main

import "github.com/spf13/viper"

func getAiAkConf() string {
	return viper.GetString("ai.ak")
}

func getAiModelConf() string {
	return viper.GetString("ai.model")
}

func getAiUrlConf() string {
	return viper.GetString("ai.url")
}
