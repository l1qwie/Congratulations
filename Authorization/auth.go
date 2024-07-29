package main

import (
	"Authorization/api/rest"
	"Authorization/app"
	"Authorization/apptype"
	resttest "Authorization/tests/rest-test"
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

func startAuthServer() {
	router := gin.Default()
	rest.AuthLogIn(router)
	rest.AuthSignIn(router)
	rest.AuthDelete(router)
	rest.AuthChange(router)
	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8090")
	err := router.RunTLS(":8090", certFile, keyFile)
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
	prepareEnv()

	go startAuthServer()
	resttest.StartAuthTest()
}
