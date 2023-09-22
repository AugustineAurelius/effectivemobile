package receiver

import (
	"context"
	"effectivemobile/FIO"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	// to consume messages
	topic := "FIO"
	partition := 0
	topicFails := "FIO_FAILED"
	connectoinForFails, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topicFails, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	connectoinForFails.SetWriteDeadline(time.Now().Add(10 * time.Second))

	batch := conn.ReadBatch(1e2, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message

	var arr []FIO.FIO

	for {
		var data map[string]interface{}

		n, err := batch.Read(b)
		if err != nil {
			break
		}

		err = json.Unmarshal(b[:n], &data)
		if err != nil {
			fmt.Errorf("%s", err)
		}

		if checkType(data) {
			fi := FIO.NewFIO(data["name"].(string),
				data["surname"].(string))

			if val, ok := data["patronymic"]; ok {
				fi.SetPatronymic(val.(string))
			}

			arr = append(arr, fi)
		} else {
			_, err = connectoinForFails.WriteMessages(
				kafka.Message{Value: []byte(b[:n])},
			)
			if err != nil {
				log.Fatal("failed to write messages:", err)
			}
		}

	}
	fmt.Println(arr)

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
	if err := connectoinForFails.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

func checkType(data map[string]interface{}) bool {
	var flg = false
	for _, v := range data {
		switch v.(type) {
		case string:
			flg = true
		default:
			return false
		}
	}
	return flg
}
