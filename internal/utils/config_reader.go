package utils

import "github.com/spf13/viper"

func ReadConfig[T any](path string, store *T) {
	viper.SetConfigName(path)      // targets jwt.yaml file
	viper.AddConfigPath("configs") // the directory where to search this file
	viper.ReadInConfig()
	viper.UnmarshalKey(path, store)
}
