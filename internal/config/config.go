package config

import "github.com/spf13/viper"

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`

	DBHost     string `mapstructure:"PG_HOST"`
	DBPort     string `mapstructure:"PG_PORT"`
	DBUser     string `mapstructure:"PG_USERNAME"`
	DBPassword string `mapstructure:"PG_PASSWORD"`
	DBName     string `mapstructure:"PG_DATABASE"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	RedisAddr string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisUser string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
