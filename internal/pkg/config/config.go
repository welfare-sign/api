package config

import "github.com/spf13/viper"

func init() {
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")
	viper.AddConfigPath("../../../config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
