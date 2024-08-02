// Модуль, в котором происходит запуск полного приложения.
// Здесь создается роутер, который передается в разные модули.
// Всё приложение выполнено в условной микросервисной архитектуре.
// Для запуска требуется добавить некоторые настройки и внести небольшие изменения,
// но в целом все модули могут работать автономно.
package main

import (
	"fmt"
	"log"

	auth "Authorization/launch"

	employees "Employees/launch"
	notif "Notifications/launch"
	sub "Subscribe/launch"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Старт сервера "Authorization"
	auth.PrepareEnv()
	auth.StartAuthServer(router)
	auth.StartAuthConsumer()

	// Старт сервера "Employees"
	employees.PrepareEnv()
	employees.StartEmployeeServer(router)
	employees.StartEmployeeConsumer()

	// Старт сервера "Notifications"
	notif.PrepareEnv()
	notif.StartNotificatorServer(router)

	// Старт сервера "Subscribe"
	sub.PrepareEnv()
	sub.StartSubscribeServer(router)
	sub.StartSubscribeConsumer()

	certFile := "keys/server.crt"
	keyFile := "keys/server.key"

	log.Print("Starting HTTPS server on :8080")
	err := router.RunTLS(":8080", certFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	}
}
