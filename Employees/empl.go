package employees

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/Congratulations/Employees/api/rest"
	resttest "github.com/l1qwie/Congratulations/Employees/tests/rest-test"
)

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

func StartEmployeeTests() {
	go StartEmployeeServer()
	resttest.StartEmployeeTests()
}
