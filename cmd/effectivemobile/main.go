package main

import (
	"context"
	"effectivemobile/FIO"
	"effectivemobile/initializers"
	"effectivemobile/kafka/producer"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"golang.org/x/crypto/sha3"
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
	config, err := initializers.LoadConfig(".") // Загружаем конфиг файл
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}
	initializers.ConnectKafka(&config) //Инициализируем читателя кафки
	initializers.ConnectDB(&config)    //Инициализируем подключение к БД
	initializers.InitializeRedis(&config)
	//Добавление методов по принципу CRUD
	router.GET("/fios", getAllFIOs)
	router.GET("/fios/:id", getFIO)
	router.POST("/fios", createFIO)
	router.DELETE("/fios/:id", deleteFIO)
	router.PUT("/fios/:id", updateFIO)
}

func main() {
	err := initializers.DB.AutoMigrate(&FIO.FIO{}) //Создаем стол путем миграции
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	fmt.Println("? Migration complete")

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{} //инициализируем вэйт группу
	wg.Add(2)               //определяем количество горутин для группы

	go func() {
		defer wg.Done()            // отложенно сообщаем о прекращении горутины
		err := router.Run(":8764") // инизиализируем сервер
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			var data map[string]interface{}               //инициализируем ассоциативный список для записини в него данных из кафки
			msg, err := initializers.KFK.ReadMessage(ctx) // читаем сообщение

			err = json.Unmarshal(msg.Value, &data)

			if flg, str := checkType(data); flg { //Проверяем корректиность наличие и корректность обязательных полей

				fi := FIO.NewFIO(data["name"].(string),
					data["surname"].(string)) // создаем новое фио

				if val, ok := data["patronymic"]; ok { // если есть отчество то, то мы его проверяем и добавляем в структуру
					if check, _ := checkValue(val.(string)); check {
						fi.SetPatronymic(val.(string))
					}
				}

				fmt.Println(fi)             //выводим в консоль
				initializers.DB.Create(&fi) // сразу кладем в БД
				sha := sha3.New256()
				key, err := sha.Write([]byte(fi.GetName() +
					fi.GetSurname() +
					fi.GetPatronymic()))
				if err != nil {
					log.Fatalln(err)
				}
				obj := new(FIO.FIO)
				err = initializers.Cache.Once(&cache.Item{
					Key:   strconv.Itoa(key),
					Value: obj, // destination
					Do: func(*cache.Item) (interface{}, error) {
						return &FIO.FIO{
							Name:       fi.GetName(),
							Surname:    fi.GetSurname(),
							Patronymic: fi.GetPatronymic(),
							Age:        fi.GetAge(),
							Gender:     fi.GetGender(),
							Nation:     fi.GetNation(),
						}, nil
					},
				})
				if err != nil {
					panic(err)
				}

			} else {

				msg.Value = append(msg.Value[:len(msg.Value)-1], []byte(`,"fail": "`+str+`"}`)...) // иначе добавляем обогащаем ошибкой
				producer.ProduceFailMessage(msg)

				if err != nil {
					log.Fatal("failed to write messages:", err)
				} // и отправляем в FIO_FAILED
			}

			if err != nil {
				panic("could not read message " + err.Error())
			}
		}
	}()

	sig := make(chan os.Signal, 1) // создаем буфферезированный канал из сигналов ос
	signal.Notify(sig, os.Interrupt)
	<-sig // ожидаем в консаль ктрл+с

	cancel()
	wg.Wait()

	err = initializers.KFK.Close() //закрываем
	if err != nil {
		log.Fatalf("Failed to close Kafka reader: %v", err)
	}

}

func updateFIO(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fio ID"})
		return
	}

	var fio FIO.FIO
	err = initializers.DB.First(&fio, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fio"})
		return
	}

	err = c.ShouldBindJSON(&fio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = initializers.DB.Save(&fio).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fio"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "FIO not found"})
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
	c.JSON(http.StatusOK, gin.H{"message": "FIO deleted"})
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
