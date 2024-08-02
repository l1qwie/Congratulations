package launch

import (
	"Authorization/api/kafka/consumer"
	"Authorization/api/rest"
	"Authorization/app"
	"Authorization/apptype"
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

func PrepareEnv() {
	var err error
	pullSymKey("keys/symmetric-key.bin")
	app.Con = new(app.Connection)
	app.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func StartAuthServer(router *gin.Engine) {
	//router := gin.Default()

	rest.AuthLogIn(router)
	rest.AuthSignIn(router)
	rest.AuthDelete(router)
	rest.AuthChange(router)

	// certFile := "keys/server.crt"
	// keyFile := "keys/server.key"

	// log.Print("Starting HTTPS server on :8090")
	// err := router.RunTLS(":8090", certFile, keyFile)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
	// }
}

func StartAuthConsumer() {
	go consumer.Consumer()
}
