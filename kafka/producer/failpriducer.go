package producer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	broker    = "localhost:9092"
	failTopic = "FIO_FAILED"
	partition = 0
)

func ProduceFailMessage(message kafka.Message) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, failTopic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		fmt.Errorf("%s", err)
	}
	message.Topic = ""
	fmt.Println(string(message.Value))

	_, err = conn.WriteMessages(message)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
