package rest

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/Congratulations/Employees/app"
	"github.com/l1qwie/Congratulations/Employees/apptype"
)

var key = []byte{}

type Err struct {
	Error error `json:"error"`
}

func decryptData(data, key []byte) ([]byte, error) {
	var (
		gcm                        cipher.AEAD
		noncesize                  int
		nonce, ciphertext, newdata []byte
	)
	block, err := aes.NewCipher(key)
	if err == nil {
		gcm, err = cipher.NewGCM(block)
	}
	if err == nil {
		noncesize := gcm.NonceSize()
		if len(data) < noncesize {
			err = fmt.Errorf("ciphertext too short")
		}
	}
	if err == nil {
		nonce, ciphertext = data[:noncesize], data[noncesize:]
	}
	if err == nil {
		newdata, err = gcm.Open(nil, nonce, ciphertext, nil)
	}
	return newdata, err
}

func UpdateEmployees(g *gin.Engine) {
	g.POST("congratulations/employees/:whatdo/:emplid", func(ctx *gin.Context) {
		var (
			statreq  int
			response interface{}
		)
		whatdo := ctx.Param("whatdo")
		diffrentemplid := ctx.Param("diffrentemplid")
		req := new(apptype.Employee)
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			statreq = http.StatusBadRequest
			response = &Err{Error: err}
		} else {
			data, err := decryptData(body, key)
			if err != nil {
				statreq = http.StatusInternalServerError
				response = &Err{Error: err}
			} else {
				err = json.Unmarshal(data, req)
				if err != nil {
					statreq = http.StatusBadRequest
					response = &Err{Error: err}
				}
			}
			if err == nil {
				answer, err := app.UpdateEmployees(req, whatdo, diffrentemplid)
				if err != nil {
					statreq = http.StatusBadRequest
					response = &Err{Error: err}
				} else {
					statreq = http.StatusOK
					response = answer
				}
			}
		}
		ctx.JSON(statreq, response)
	})
}

func GetEmployees(g *gin.Engine) {
	g.GET("congratulations/employees/:id/:limit", func(ctx *gin.Context) {
		log.Print("Someone has just called congratulations/employees/:id/:limit")
		var (
			statreq, id, limit int
			response           interface{}
			err                error
		)
		idstr := ctx.Param("id")
		limitstr := ctx.Param("limit")
		id, err = strconv.Atoi(idstr)
		if err != nil {
			log.Print("Something went wrong when tryed to convert string to int in the first time")
			statreq = http.StatusBadRequest
			response = &Err{Error: err}
		} else {
			limit, err = strconv.Atoi(limitstr)
			if err != nil {
				log.Print("Something went wrong when tryed to convert string to int in the second time")
				statreq = http.StatusBadRequest
				response = &Err{Error: err}
			}
		}
		if err == nil {
			log.Print("Everything was ok. Starts the bussines logic")
			employees, err := app.GetEmployees(id, limit)
			if err != nil {
				log.Printf("Something went wrong during the bussines logic: %s", err)
				statreq = http.StatusBadRequest
				response = &Err{Error: err}
			} else {
				log.Print("Everything was ok. Starting sending the answer")
				statreq = http.StatusOK
				response = employees
			}
		}
		ctx.JSON(statreq, response)
		log.Print("Server congratulations/employees/:id/:limit has just fihished its job")
	})
}
