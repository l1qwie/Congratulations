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

	"Employees/app"
	"Employees/apptype"

	"github.com/gin-gonic/gin"
)

type Err struct {
	Error string `json:"error"`
}

type Inside struct {
	Whatdo     string `form:"whatdo"`
	EmployeeId string `form:"emplid"`
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

// Обрабатывает запрос с post endPoint'а
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
			err = json.Unmarshal(data, empl)
			if err != nil {
				log.Print("Couldn't unmarshal data to req (apptype.Employee)")
				statreq = http.StatusBadRequest
				respErr.Error = err.Error()
			}
		}
	}
	return statreq, err
}

// Обрабатывает запрос с get endPoint'а
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

// @Summary Update any kind of information about an employee
// @Description Update or modify employee information based on the action specified in {whatdo}.
// The {whatdo} parameter can be "new", "update", "delete", "sub", or "unsub". Depending on the action, the request
// payload and the parameters required will vary:
// - "new": Adds a new employee. Do not provide {emplid}. Send the employee data as JSON.
// - "update": Updates an existing employee. Provide all employee details in the JSON payload and specify {emplid}.
// - "delete": Deletes an employee. Provide only the Id field in the JSON payload.
// - "sub": Subscribes one employee to another. Provide the Id field in the JSON payload for the employee being subscribed and {emplid} for the target employee.
// - "unsub": Unsubscribes one employee from another. Provide the Id field in the JSON payload for the employee being unsubscribed and {emplid} for the target employee.
// If the process is successful, the response will be a string "The employee has been updated".
// In case of an error, the response will be a JSON object of rest.Err. If there is an encryption error, a plain string will be returned.
// @Accept json
// @Produce string
// @Param   whatdo      path    string  			false   "What app should do: new, update, delete, sub, unsub"
// @Param   emplid	    path    string  			false   "The second Employee ID"
// @Param 	employee 	body 	apptype.Employee 	true 	"Employee details"
// @Success 200 {string}  	string "The employee has been updated"
// @Failure 400 {object} 	rest.Err
// @Failure 400 {string} 	string "Something went wrong while we were trying to encrypt data to send"
// @Router /congratulations/employees/{whatdo}/{emplid} [post]
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

// @Summary Get Employees Info
// @Description Retrieves an array of employees based on the specified employee ID and limit.
// If you want to get an array with an employee's subscription, you must specify the ID.
// To get a variable length of array, specify the limit. The default limit is 15.
// If successful, the response will be an array of *apptype.Employee.
// In case of an error, the response will be a JSON object of rest.Err.
// If the response is a plain string, it indicates an encryption error. Please contact the developer.
// @Accept url
// @Produce json
// @Param   id      path    string  false   "Employee ID"
// @Param   limit   path    string  false   "Limit"
// @Success 200 {array}  apptype.Employee
// @Failure 400 {object} rest.Err
// @Failure 400 {string} string "Something went wrong while we were trying to encrypt data to send"
// @Router /congratulations/employees/{id}/{limit} [get]
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
