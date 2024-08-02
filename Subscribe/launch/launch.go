package launch

import (
	"Subscribe/api/kafka/consumer"
	"Subscribe/api/rest"
	"Subscribe/app"
	"Subscribe/apptype"
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

func StartSubscribeServer(router *gin.Engine) {
	rest.SubscribeServer(router)
}
func StartSubscribeConsumer() {
	go consumer.Consumer()
}
