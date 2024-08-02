package main

import (
	"fmt"
	"log"

	"Subscribe/api/kafka/consumer"
	"Subscribe/api/rest"
	"Subscribe/apptype"
	"Subscribe/launch"
	kafkatest "Subscribe/tests/kafka-test"
	resttest "Subscribe/tests/rest-test"

	"github.com/gin-gonic/gin"
)

func prepareTestEnv() {
	var err error
	//pullSymKey("keys/symmetric-key.bin")
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
	launch.PrepareEnv()
	prepareTestEnv()

	go StartSubscribeServer()
	resttest.StartSubTest()

	go consumer.Consumer()
	kafkatest.StartSubKafkaTest()
}
