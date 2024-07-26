package main

import (
	"Notifications/api/rest"
	"Notifications/app"
	"Notifications/apptype"
	"Notifications/tests"
	"fmt"
	"log"
	"os"

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

func prepareEnv() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	app.Con = new(app.Connection)
	app.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Print("Servese Notifications has been launched")

	prepareEnv()

	go startNotifierServer()
	go app.Guardian()

	tests.TestNotifications()
	tests.UnitTest()

	log.Print("Servese Notifications has finished its work")
}
