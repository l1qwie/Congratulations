package launch

import (
	"Employees/api/kafka/consumer"
	"Employees/api/rest"
	"Employees/app"
	"Employees/apptype"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Добавляет редис-клиента
func addClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "congratulations-redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	return client, err
}

// Вынимает данные о ключе шифрования из файла (сгенерирован до запуска приложения)
func pullSymKey(filePath string) {
	var err error
	apptype.SymKey, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
}

func PrepareEnv() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	app.Client = new(app.RedClient)
	app.Client.Cl, err = addClient()
	if err != nil {
		panic(err)
	}
	log.Print("The connection to redis is successful")
}

func StartEmployeeServer(router *gin.Engine) {
	rest.GetEmployees(router)
	rest.UpdateEmployees(router)
}

func StartEmployeeConsumer() {
	go consumer.Consumer()
}
