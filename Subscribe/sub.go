package main

import (
	"fmt"
	"log"
	"os"

	"Subscribe/api/kafka/consumer"
	"Subscribe/api/rest"
	"Subscribe/app"
	"Subscribe/apptype"
	kafkatest "Subscribe/tests/kafka-test"
	resttest "Subscribe/tests/rest-test"

	"github.com/gin-gonic/gin"
)

// Вынимает данные о ключе шифрования из файла (сгенерирован до запуска приложения)
func pullSymKey(filePath string) {
	var err error
	apptype.SymKey, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
}

func prepareEnv() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	app.Con = new(app.Connection)
	app.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func prepareTestEnv() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	resttest.Con = new(resttest.TestConnection)
	resttest.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func StartSubscribeServer() {
	router := gin.Default()
	rest.SubscribeServer(router)
	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8019")
	err := router.RunTLS(":8019", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}

func main() {
	prepareEnv()
	prepareTestEnv()

	// go StartSubscribeServer()
	// resttest.StartSubTest()

	go consumer.Consumer()
	kafkatest.StartSubKafkaTest()
}
