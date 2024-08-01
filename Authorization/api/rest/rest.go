package rest

import (
	"Authorization/app"
	"Authorization/apptype"
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

func handleRequest(reqbody io.ReadCloser, respErr *Error) ([]byte, int) {
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
		}
	}
	return data, statreq
}

// Обрабатывает запрос
func handleRequestLogSign(reqbody io.ReadCloser, auth *apptype.Auth, respErr *Error) (int, error) {
	var (
		statreq int
		data    []byte
		err     error
	)
	data, statreq = handleRequest(reqbody, respErr)
	if respErr.Err == "" {
		log.Printf("Successfully decrypted data: %x", data)
		err = json.Unmarshal(data, auth)
		if err != nil {
			log.Print("Couldn't unmarshal data to req (apptype.Auth)")
			statreq = http.StatusBadRequest
			respErr.Err = err.Error()
		}
	}
	return statreq, err
}

func handleRequestDel(reqbody io.ReadCloser, del *apptype.Delete, respErr *Error) (int, error) {
	var err error
	data, statreq := handleRequest(reqbody, respErr)
	if respErr.Err == "" {
		log.Printf("Successfully decrypted data: %x", data)
		err = json.Unmarshal(data, del)
		if err != nil {
			log.Print("Couldn't unmarshal data to req (apptype.Delete)")
			statreq = http.StatusBadRequest
			respErr.Err = err.Error()
		}
	}
	return statreq, err
}

func handleRequestChange(reqbody io.ReadCloser, changes *apptype.Change, respErr *Error) (int, error) {
	var err error
	data, statreq := handleRequest(reqbody, respErr)
	if respErr.Err == "" {
		log.Printf("Successfully decrypted data: %x", data)
		err = json.Unmarshal(data, changes)
		if err != nil {
			log.Print("Couldn't unmarshal data to req (apptype.Change)")
			statreq = http.StatusBadRequest
			respErr.Err = err.Error()
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

// @Summary Log In
// @Description If you want to log in you must send your nickname and password
// @Accept json
// @Produce string
// @Param	auth	body	apptype.Auth	true	"Details for logging in"
// @Success 200 {string}  string "You have been authrizated"
// @Failure 400 {object} rest.Error
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/authorization/login" [post]
func AuthLogIn(router *gin.Engine) {
	router.POST("congratulations/authorization/login", func(ctx *gin.Context) {
		log.Print("Someone has just called the server congratulations/authorization/login")
		var (
			statreq    int
			bodyanswer []byte
			answer     string
			err        error
		)
		auth := new(apptype.Auth)
		respError := new(Error)
		statreq, err = handleRequestLogSign(ctx.Request.Body, auth, respError)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			answer, err = app.LogIn(auth, ctx.ClientIP())
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
		log.Print("Server congratulations/authorization/login has just finished its job")
	})
}

// @Summary Sign up
// @Description If you want to sign up you must send your nickname and password and it will put this information to the database.
// The endpoint will return you a json-object of your information
// @Accept json
// @Produce json
// @Param	notified	body	apptype.Auth	true	"Details for signing up"
// @Success 200 {object} apptype.SignIn
// @Failure 400 {object} rest.Error
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/authorization/signin [post]
func AuthSignIn(router *gin.Engine) {
	router.POST("congratulations/authorization/signin", func(ctx *gin.Context) {
		log.Print("Someone has just called the server congratulations/authorization/signin")
		var (
			statreq    int
			bodyanswer []byte
			err        error
		)
		auth := new(apptype.Auth)
		signin := new(apptype.SignIn)
		respError := new(Error)
		statreq, err = handleRequestLogSign(ctx.Request.Body, auth, respError)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			signin, err = app.SignUp(auth, ctx.ClientIP())
			if err != nil {
				statreq = http.StatusBadRequest
				log.Printf("Something went wrong during the business logic: %s", err)
				respError.Err = err.Error()
			} else {
				log.Printf("signin from business logic: %v", signin)
				statreq = http.StatusOK
			}
		}
		if statreq != http.StatusOK {
			bodyanswer, err = json.Marshal(respError)
			if err != nil {
				log.Printf("Couldn't marshal RespErr. Something went wrong: %s", err)
			}
		} else {
			bodyanswer, err = json.Marshal(signin)
			if err != nil {
				log.Printf("Couldn't marshal signin. Something went wrong : %s", err)
			}
		}

		response := prepareResponse(bodyanswer)
		ctx.Data(statreq, "application/octet-stream", response)
		log.Print("Server congratulations/authorization/signin has just finished its job")
	})
}

// @Summary Delete data about you
// @Description If you want to delete all data about you in the database you should call this endpoint.
// @Accept json
// @Produce string
// @Param	delete		body		apptype.Delete		true	"Details for deleting"
// @Success 200 {string} string "Employee has been deleted"
// @Failure 400 {object} rest.Error
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/authorization/delete [delete]
func AuthDelete(router *gin.Engine) {
	router.DELETE("congratulations/authorization/delete", func(ctx *gin.Context) {
		log.Print("Someone has just called the server congratulations/authorization/delete")
		var (
			statreq    int
			bodyanswer []byte
			err        error
			answer     string
		)
		del := new(apptype.Delete)
		respError := new(Error)
		statreq, err = handleRequestDel(ctx.Request.Body, del, respError)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			answer, err = app.DeleteEmployee(del, ctx.ClientIP())
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
		log.Print("Server congratulations/authorization/delete has just finished its job")
	})
}

// @Summary Change data about you
// @Description If you want to change some data about you in the database, you should call this endpoint.
// You can change only you nickname or password.
// @Accept json
// @Produce string
// @Param	change		body		apptype.Change		true	"Details for changing"
// @Success 200 {string} string "Employee has been deleted"
// @Failure 400 {object} rest.Error
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to sends"
// @Router /congratulations/authorization/delete [delete]
func AuthChange(router *gin.Engine) {
	router.PUT("congratulations/authorization/change", func(ctx *gin.Context) {
		log.Print("Someone has just called the server congratulations/authorization/change")
		var (
			statreq    int
			bodyanswer []byte
			err        error
			answer     string
		)
		changes := new(apptype.Change)
		respError := new(Error)
		statreq, err = handleRequestChange(ctx.Request.Body, changes, respError)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			answer, err = app.ChangeEmployee(changes, ctx.ClientIP())
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
		log.Print("Server congratulations/authorization/change has just finished its job")
	})
}
