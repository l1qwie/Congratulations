package rest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

type Err struct {
	Error string `json:"error"`
}

type Inside struct {
	Whatdo     string `form:"whatdo"`
	EmployeeId string `form:"emplid"`
}

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

func handlePostRequest(ctx *gin.Context, in *Inside, empl *apptype.Employee, respErr *Err) (int, error) {
	var (
		statreq    int
		data, body []byte
		err        error
	)
	in.Whatdo = ctx.Param("whatdo")
	in.EmployeeId = ctx.Param("emplid")
	log.Printf("in.WhatDo: %s, in.EmployeeId: %s", in.Whatdo, in.EmployeeId)
	body, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Print("Couldn't read from ctx.Request.Body")
		statreq = http.StatusBadRequest
		respErr.Error = err.Error()
	} else {
		log.Printf("Successfully read from ctx.Request.Body: %x", body)
		data, err = decryptData(body, apptype.SymKey)
		if err != nil {
			log.Printf("Couldn't decrypt data from ctx.Request.Body (body): %s", err)
			statreq = http.StatusInternalServerError
			respErr.Error = err.Error()
		} else {
			log.Printf("Successfully decrypted data: %x", data)
			if in.Whatdo != "delete" {
				err = json.Unmarshal(data, empl)
				if err != nil {
					log.Print("Couldn't unmarshal data to req (apptype.Employee)")
					statreq = http.StatusBadRequest
					respErr.Error = err.Error()
				}
			}
		}
	}
	return statreq, err
}

func handleGetRequest(ctx *gin.Context, respErr *Err) (int, int, int, error) {
	var (
		id, limit, statreq int
		err                error
	)
	idstr := ctx.Param("id")
	limitstr := ctx.Param("limit")
	id, err = strconv.Atoi(idstr)
	if err != nil {
		log.Print("Something went wrong when tried to convert string to int in the first time")
		statreq = http.StatusBadRequest
		respErr.Error = err.Error()
	} else {
		limit, err = strconv.Atoi(limitstr)
		if err != nil {
			log.Print("Something went wrong when tried to convert string to int in the second time")
			statreq = http.StatusBadRequest
			respErr.Error = err.Error()
		}
	}
	return id, limit, statreq, err
}

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

func UpdateEmployees(g *gin.Engine) {
	g.POST("congratulations/employees/:whatdo/:emplid", func(ctx *gin.Context) {
		log.Print("Someone has just called congratulations/employees/:whatdo/:emplid")
		var (
			statreq    int
			bodyanswer []byte
			answer     string
			err        error
		)
		in := new(Inside)
		req := new(apptype.Employee)
		respErr := new(Err)
		statreq, err = handlePostRequest(ctx, in, req, respErr)
		if err == nil {
			log.Print("Everything was ok. Started the business logic")
			answer, err = app.UpdateEmployees(req, in.Whatdo, in.EmployeeId)
			if err != nil {
				log.Printf("Something went wrong during the business logic: %s", err)
				statreq = http.StatusBadRequest
				respErr.Error = err.Error()
			} else {
				statreq = http.StatusOK
			}
		}
		if statreq != http.StatusOK {
			bodyanswer, err = json.Marshal(respErr)
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

func GetEmployees(g *gin.Engine) {
	g.GET("congratulations/employees/:id/:limit", func(ctx *gin.Context) {
		log.Print("Someone has just called congratulations/employees/:id/:limit")
		var (
			statreq, id, limit int
			err                error
			bodyanswer         []byte
		)
		respErr := new(Err)
		id, limit, statreq, err = handleGetRequest(ctx, respErr)
		if err == nil {
			log.Print("Everything was ok. Starts the business logic")
			employees, err := app.GetEmployees(id, limit)
			if err != nil {
				log.Printf("Something went wrong during the business logic: %s", err)
				statreq = http.StatusBadRequest
				respErr.Error = err.Error()
			} else {
				bodyanswer, err = json.Marshal(employees)
				if err != nil {
					log.Printf("Couldn't marshal array employees: %s", err)
					statreq = http.StatusBadRequest
					respErr.Error = err.Error()
				} else {
					log.Print("Successfully marshaled a slice of employees")
					statreq = http.StatusOK
				}
			}
		}

		if statreq != http.StatusOK {
			bodyanswer, err = json.Marshal(respErr)
			if err != nil {
				log.Printf("Couldn't marshal RespErr. Something went wrong: %s", err)
			}
		}

		response := prepareResponse(bodyanswer)
		ctx.Data(statreq, "application/octet-stream", response)
		log.Print("Server congratulations/employees/:id/:limit has just finished its job")
	})
}
