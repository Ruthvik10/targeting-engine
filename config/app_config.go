package config

import "github.com/spf13/viper"

type AppConfig struct {
	DBURI          string `mapstructure:"MONGO_URI"`
	ServerAddr     string `mapstructure:"SERVER_ADDR"`
	RedisURI       string `mapstructure:"REDIS_URI"`
	RedisCacheExpr int    `mapstructure:"REDIS_CACHE_EXPR"`
	TestDBURI      string `mapstructure:"TEST_MONGO_URI"`
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
