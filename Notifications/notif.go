package main

import (
	"Notifications/api/rest"
	"Notifications/app"
	"Notifications/launch"
	"Notifications/tests"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func startNotifierServer() {
	router := gin.Default()
	rest.Notify(router)
	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8099")
	err := router.RunTLS(":8089", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}

func main() {
	log.Print("Servese Notifications has been launched")

	launch.PrepareEnv()

	go startNotifierServer()
	go app.Guardian()

	tests.TestNotifications()
	tests.UnitTest()

	log.Print("Servese Notifications has finished its work")
}
