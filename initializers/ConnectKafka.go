package initializers

import (
	"fmt"
	"github.com/segmentio/kafka-go"
)

var KFK *kafka.Reader

func ConnectKafka(config *Config) {

	KFK = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.BrokerKafka},
		Topic:   config.Topic,
	})
	fmt.Println("? Connected Successfully to the Kafka for listen")
}
