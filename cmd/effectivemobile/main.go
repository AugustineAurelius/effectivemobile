package main

import (
	"context"
	"effectivemobile/FIO"
	"effectivemobile/initializers"
	"effectivemobile/kafka/producer"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
)

var router = gin.Default()

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}
	initializers.ConnectKafka(&config)
	initializers.ConnectDB(&config)

	router.GET("/fios", getAllFIOs)

	// Обработка GET-запроса для получения пользователя по ID
	router.GET("/fios/:id", getFIO)

	// Обработка POST-запроса для создания пользователя
	router.POST("/fios", createFIO)

	// Обработка DELETE-запроса для удаления пользователя по ID
	router.DELETE("/users/:id", deleteFIO)
	router.PUT("/users/:id", updateFIO)
}

func main() {
	initializers.DB.AutoMigrate(&FIO.FIO{})
	fmt.Println("? Migration complete")

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := router.Run(":8764")
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			var data map[string]interface{}
			msg, err := initializers.KFK.ReadMessage(ctx)

			err = json.Unmarshal(msg.Value, &data)

			if flg, str := checkType(data); flg {

				fi := FIO.NewFIO(data["name"].(string),
					data["surname"].(string))

				if val, ok := data["patronymic"]; ok {
					fi.SetPatronymic(val.(string))
				}

				fmt.Println(fi)
				initializers.DB.Create(&fi)

			} else {
				msg.Value = append(msg.Value[:len(msg.Value)-1], []byte(`,"fail": "`+str+`"}`)...)
				go producer.ProduceFailMessage(msg)
			}

			if err != nil {
				panic("could not read message " + err.Error())
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	cancel()
	wg.Wait()

	err := initializers.KFK.Close()
	if err != nil {
		log.Fatalf("Failed to close Kafka reader: %v", err)
	}

}

func updateFIO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var fio FIO.FIO
	err = initializers.DB.First(&fio, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	err = c.ShouldBindJSON(&fio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = initializers.DB.Save(&fio).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, fio)
}
func getAllFIOs(c *gin.Context) {
	var fios []FIO.FIO
	err := initializers.DB.Find(&fios).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fios)
}
func getFIO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user FIO.FIO
	err = initializers.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func createFIO(c *gin.Context) {
	var user FIO.FIO
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = initializers.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
func deleteFIO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = initializers.DB.Delete(&FIO.FIO{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
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
