package main

import (
	"context"
	"effectivemobile/FIO"
	"effectivemobile/kafka/producer"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"regexp"
)

const (
	broker = "localhost:9092"
	topic  = "FIO"
)

func main() {
	consume(context.Background())
}

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		var data map[string]interface{}
		msg, err := r.ReadMessage(ctx)

		err = json.Unmarshal(msg.Value, &data)

		if flg, str := checkType(data); flg {
			fi := FIO.NewFIO(data["name"].(string),
				data["surname"].(string))

			if val, ok := data["patronymic"]; ok {
				fi.SetPatronymic(val.(string))
			}

		} else {
			fmt.Println("err", str)
			msg.Value = append(msg.Value[:len(msg.Value)-1], []byte(`,"fail": "`+str+`"}`)...)
			go producer.ProduceFailMessage(msg)
		}

		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func checkType(data map[string]interface{}) (bool, string) {
	var keys []string

	var flg = false
	for k, v := range data {
		switch v.(type) {
		case string:
			keys = append(keys, k)
			flg = true
		default:
			return false, "Неверный тип значений"
		}
	}
	_, str := checkKeys(keys, "name")
	if str != "" {
		return false, str
	}
	_, str = checkKeys(keys, "surname")
	if str != "" {
		return false, str
	}

	_, str = checkValue(data["name"].(string))
	if str != "" {
		return false, str
	}
	_, str = checkValue(data["surname"].(string))
	if str != "" {
		return false, str
	}

	return flg, ""
}

func checkKeys(s []string, str string) (bool, string) {
	for _, v := range s {
		if v == str {
			return true, ""
		}
	}
	return false, "Нет обязательного поля"
}

func checkValue(str string) (bool, string) {
	re := regexp.MustCompile("[0-9]+")
	if 0 != len(re.FindAllString(str, -1)) {
		return false, "Некоректные значения"
	}
	return true, ""
}
