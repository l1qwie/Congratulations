package main

import (
	"fmt"
	"log"
	"os"

	"Employees/api/rest"
	"Employees/apptype"

	resttest "Employees/tests/rest-test"

	"github.com/gin-gonic/gin"
)

// Запускает микросервис Employees
func StartEmployeeServer() {
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
	pullSymKey("keys/symmetric-key.bin")
	go StartEmployeeServer()
	resttest.StartEmployeeTests()
}
