package rest

import (
	"Notifications/app"
	"Notifications/apptype"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Err string `json:"error"`
}

// Расшифровывет данные с помощью ключа
func decryptData(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	if len(nonce) != nonceSize {
		return nil, fmt.Errorf("nonce size is incorrect")
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
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

// Обрабатывает запрос
func handleRequest(reqbody io.ReadCloser, notif *apptype.Notified, respErr *Error) (int, error) {
	var (
		statreq    int
		data, body []byte
		err        error
	)
	body, err = io.ReadAll(reqbody)
	if err != nil {
		log.Print("Couldn't read from ctx.Request.Body")
		statreq = http.StatusBadRequest
		respErr.Err = err.Error()
	} else {
		log.Printf("Successfully read from ctx.Request.Body: %x", body)
		data, err = decryptData(body, apptype.SymKey)
		if err != nil {
			log.Printf("Couldn't decrypt data from ctx.Request.Body (body): %s", err)
			statreq = http.StatusInternalServerError
			respErr.Err = err.Error()
		} else {
			log.Printf("Successfully decrypted data: %x", data)
			err = json.Unmarshal(data, notif)
			if err != nil {
				log.Print("Couldn't unmarshal data to req (apptype.Employee)")
				statreq = http.StatusBadRequest
				respErr.Err = err.Error()
			}
		}
	}
	return statreq, err
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

// @Summary This API can notify any person by sending email
// @Description You must give all fields in structer and then the server can send email by using "gomail"
// @Accept json
// @Produce string
// @Param	notified	body	apptype.Notified	true	"Details about the subscriber and about 'birthday-boy'"
// @Success 200 {string}  string "The employee was notified"
// @Failure 400 {object} rest.Err
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/notify [post]
func Notify(router *gin.Engine) {
	router.POST("/congratulations/notify", func(ctx *gin.Context) {
		var (
			statreq    int
			bodyanswer []byte
			answer     string
			err        error
		)
		notified := new(apptype.Notified)
		respError := new(Error)
		statreq, err = handleRequest(ctx.Request.Body, notified, respError)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			answer, err = app.Notify(notified)
			if err != nil {
				statreq = http.StatusBadRequest
				log.Printf("Something went wrong during the business logic: %s", err)
				respError.Err = err.Error()
			} else {
				statreq = http.StatusOK
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
		log.Print("Server congratulations/employees/:whatdo/:emplid has just finished its job")
	})
}
