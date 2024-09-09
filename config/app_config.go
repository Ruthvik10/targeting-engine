package config

import "github.com/spf13/viper"

type AppConfig struct {
	DBURI      string `mapstructure:"MONGO_URI"`
	ServerAddr string `mapstructure:"SERVER_ADDR"`
}

func Load(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
