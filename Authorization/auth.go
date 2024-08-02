package main

import (
	"Authorization/api/kafka/consumer"
	"Authorization/apptype"
	"Authorization/launch"
	kafkatest "Authorization/tests/kafka-test"
	resttest "Authorization/tests/rest-test"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func startAuthServer() {
	router := gin.Default()
	launch.StartAuthServer(router)
	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8090")
	err := router.RunTLS(":8090", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}

func prepareTestEnv() {
	var err error
	resttest.Con = new(resttest.TestConnection)
	resttest.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func main() {
	launch.PrepareEnv()
	prepareTestEnv()

	go startAuthServer()
	resttest.StartAuthTest()

	go consumer.Consumer()
	kafkatest.StartTestConsumer()
}
