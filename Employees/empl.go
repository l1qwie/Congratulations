package employees

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/Congratulations/Employees/api/rest"
	resttest "github.com/l1qwie/Congratulations/Employees/tests/rest-test"
)

// Запускает микросервис Employees
func StartEmployeeServer() {
	router := gin.Default()
	rest.GetEmployees(router)
	rest.UpdateEmployees(router)
	certFile := "server.crt"
	keyFile := "server.key"

	log.Print("Starting HTTPS server on :8099")
	err := router.RunTLS(":8099", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}

// Запускает сервер и тесты для микросервиса Employees
func StartEmployeeTests() {
	go StartEmployeeServer()
	resttest.StartEmployeeTests()
}
