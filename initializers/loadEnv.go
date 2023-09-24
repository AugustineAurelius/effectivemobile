package initializers

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	BrokerKafka     string `mapstructure:"BROKER_KAFKA"`
	FailBrokerKafka string `mapstructure:"FAIL_BROKER_KAFKA"`
	Topic           string `mapstructure:"TOPIC"`
	FailTopic       string `mapstructure:"FAIL_TOPIC"`
	PartitionKafka  int    `mapstructure:"PARTITION_KAFKA"`

	RedisServer1 string `mapstructure:"REDIS_SERVER_1"`
	RedisServer2 string `mapstructure:"REDIS_SERVER_2"`
}

func LoadConfig(path string) (config Config, err error) {
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
