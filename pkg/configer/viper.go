package configer

import (
	"log"

	"github.com/spf13/viper"
)

type ServerData struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	ApiPort string `mapstructure:"api-port"`
}

type PostgresqlData struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db-name"`
	Port     string `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Sslmode  string `mapstructure:"sslmode"`
}

type Config struct {
	Server     ServerData     `mapstructure:"server-data"`
	Postgresql PostgresqlData `mapstructure:"postgresql-data"`
}

var AppConfig Config

func InitConfig(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal(err)
	}
}
