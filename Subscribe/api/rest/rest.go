package rest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"Subscribe/app"
	"Subscribe/apptype"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Err string `json:"error"`
}

// Шифрует данные с помощью ключа
func encryptData(data, key []byte) ([]byte, error) {
	var (
		gcm   cipher.AEAD
		nonce []byte
	)
	block, err := aes.NewCipher(key)
	if err == nil {
		gcm, err = cipher.NewGCM(block)
	}
	if err == nil {
		nonce = make([]byte, gcm.NonceSize())
		_, err = io.ReadFull(rand.Reader, nonce)
	}
	return gcm.Seal(nonce, nonce, data, nil), err
}

// Подготавливает ответ от сервера клиенту
func prepareResponse(body []byte) []byte {
	var (
		response, encryptMes []byte
		statreq              int
		err                  error
	)
	log.Printf("Body to encode: %s", string(body))
	encryptMes, err = encryptData(body, apptype.SymKey)
	if err != nil {
		log.Printf("Couldn't encrypt message answer. Something went wrong: %s", err)
		statreq = http.StatusBadRequest
		response = []byte("Something went wrong while we were trying to encrypt data to send")
	} else {
		statreq = http.StatusOK
		response = encryptMes
	}
	log.Print(statreq, response)
	return response
}

func handeWholeRequest(ctx *gin.Context, whatdo string) {
	var (
		statreq, subid, subto int
		answer                string
		err                   error
		bodyanswer            []byte
	)
	subscriberstr := ctx.Param("subscriber")
	subtostr := ctx.Param("subto")
	respError := new(Error)
	subid, err = strconv.Atoi(subscriberstr)
	if err != nil {
		statreq = http.StatusBadRequest
		respError.Err = err.Error()
	} else {
		subto, err = strconv.Atoi(subtostr)
		if err != nil {
			statreq = http.StatusBadRequest
			respError.Err = err.Error()
		} else {
			answer, err = app.SubscribeAction(subid, subto, whatdo)
			if err != nil {
				statreq = http.StatusBadRequest
				respError.Err = err.Error()
			} else {
				statreq = http.StatusOK
			}
		}
	}
	if statreq != http.StatusOK {
		bodyanswer, err = json.Marshal(respError)
		if err != nil {
			log.Printf("Couldn't marshal RespErr. Something went wrong: %s", err)
		}
	} else {
		bodyanswer = []byte(answer)
	}

	response := prepareResponse(bodyanswer)
	ctx.Data(statreq, "application/octet-stream", response)
}

// @Summary This API can subscribe someone to someone.
// @Description There's only one function - subscribe
// @Accept path
// @Produce string
// @Param	SubscriberId	path	integer		true	"Employee's id who's subscribed to somemone (SubtoId)"
// @Param 	SubtoId			path	integer 	true	"Employee's id who's subscriber (SubscriberId)"
// @Success 200 {string}  string "You've just subscribed to {id}"
// @Failure 400 {object} rest.Err
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/subscriptions/sub/{subscriber}/{subto} [put]
func sub(router *gin.Engine) {
	router.PUT("congratulations/subscriptions/sub/:subscriber/:subto", func(ctx *gin.Context) {
		log.Print("Someone has just called ongratulations/subscriptions/sub/:subscriber/:subto server")
		handeWholeRequest(ctx, app.Sub)
		log.Print("Server ongratulations/subscriptions/sub/:subscriber/:subto has just finished its job")
	})
}

// @Summary This API can unsubscribe someone from someone.
// @Description There's only one functions - unsubscribe
// @Accept path
// @Produce string
// @Param	SubscriberId	path	integer		true	"Employee's id who's subscribed to somemone (SubtoId)"
// @Param 	SubtoId			path	integer 	true	"Employee's id who's subscriber (SubscriberId)"
// @Success 200 {string}  string "You've just unsubscribed from {id}"
// @Failure 400 {object} rest.Err
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/subscriptions/unsub/{subscriber}/{subto} [put]
func unsub(router *gin.Engine) {
	router.PUT("congratulations/subscriptions/unsub/:subscriber/:subto", func(ctx *gin.Context) {
		log.Print("Someone has just called ongratulations/subscriptions/unsub/:subscriber/:subto server")
		handeWholeRequest(ctx, app.Unsub)
		log.Print("Server ongratulations/subscriptions/unsub/:subscriber/:subto has just finished its job")
	})
}

func SubscribeServer(router *gin.Engine) {
	sub(router)
	unsub(router)
}
