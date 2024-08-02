package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"Employees/api/kafka/consumer"
	"Employees/api/rest"
	"Employees/app"
	"Employees/apptype"

	"Employees/tests/kafka"
	resttest "Employees/tests/rest-test"

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

// Запускает микросервис Employees
func startEmployeeServer() {
	router := gin.Default()
	rest.GetEmployees(router)
	rest.UpdateEmployees(router)
	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8099")
	err := router.RunTLS(":8099", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}

// Вынимает данные о ключе шифрования из файла (сгенерирован до запуска приложения)
func pullSymKey(filePath string) {
	var err error
	apptype.SymKey, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
}

// Запускает сервер и тесты для микросервиса Employees
func main() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	app.Client = new(app.RedClient)
	app.Client.Cl, err = addClient()
	if err != nil {
		panic(err)
	}
	log.Print("The connection to redis is successful")

	go startEmployeeServer()
	go consumer.Consumer()

	resttest.StartEmployeeTests()
	kafka.TestConsumer()
}
