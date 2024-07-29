package resttest

import (
	"Authorization/apptype"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Создает клиента для https
func createClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // only for testing
			},
		},
	}
	return client
}

// Шифрует данные отправляемые на сервер
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

// Расшифровка данных полученных от сервера
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

// func deleteAuth()

func deleteAuthDelete(body []byte) string {
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("DELETE", "https://localhost:8090/congratulations/authorization/delete", bytes.NewBuffer(encryptedBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to do request: %s", err))
	}

	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}

	log.Printf("Response: %v", respbody)

	decryptedMessage, err := decryptData(respbody, apptype.SymKey)
	if err != nil {
		log.Print(string(respbody))
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}
	res := string(decryptedMessage)
	log.Print("Response:", res)
	return res
}

func postAuthSignIn(body []byte) *apptype.SignIn {
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("POST", "https://localhost:8090/congratulations/authorization/signin", bytes.NewBuffer(encryptedBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to do request: %s", err))
	}

	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}

	log.Printf("Response: %v", respbody)

	decryptedMessage, err := decryptData(respbody, apptype.SymKey)
	if err != nil {
		log.Print(string(respbody))
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}
	res := new(apptype.SignIn)
	err = json.Unmarshal(decryptedMessage, res)
	if err != nil {
		panic(err)
	}
	log.Print("Response:", res)
	return res
}

func postAuthLogIn(body []byte) string {
	var res string
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("POST", "https://localhost:8090/congratulations/authorization/login", bytes.NewBuffer(encryptedBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to do request: %s", err))
	}

	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}

	log.Printf("Response: %v", respbody)

	decryptedMessage, err := decryptData(respbody, apptype.SymKey)
	if err != nil {
		log.Print(string(respbody))
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}

	res = string(decryptedMessage)
	log.Print("Response:", res)
	return res
}

func putAuthChange(body []byte) string {
	var res string
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("PUT", "https://localhost:8090/congratulations/authorization/change", bytes.NewBuffer(encryptedBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed to do request: %s", err))
	}

	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}

	log.Printf("Response: %v", respbody)

	decryptedMessage, err := decryptData(respbody, apptype.SymKey)
	if err != nil {
		log.Print(string(respbody))
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}

	res = string(decryptedMessage)
	log.Print("Response:", res)
	return res
}

func testLogIn() {
	defer Con.resetSequence()
	defer Con.deleteEmployees()

	Con.createEmployee()

	log.Print("Test testLogIn() has been started")
	auth := &apptype.Auth{
		Nickname: "mama-miya",
		Password: "12345678",
	}
	body, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}
	answer := postAuthLogIn(body)
	if answer != "You have been authrizated" {
		panic(fmt.Sprintf("EXPECTED: You have been authrizated, RESIEVED: %s", answer))
	}
	if !Con.checkLoggedIn(auth) {
		panic("Client didn't log in")
	}
	log.Print("Test testLogIn has been finished")
}

func testSingIn() {
	defer Con.resetSequence()
	defer Con.deleteEmployees()

	Con.createEmployee()

	log.Print("Test testSignIn has been starterd")
	auth := &apptype.Auth{
		Nickname: "pumpalumpa",
		Password: "1029384756",
	}
	body, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}
	resp := new(apptype.SignIn)
	resp = postAuthSignIn(body)
	if resp.Answer != "Welcome! You have successfuly sign in!" {
		panic(fmt.Sprintf(`EXPECTED: resp.Answer = "Welcome! You have successfuly sign in!". RECIEVED: resp.Answer = %s`, resp.Answer))
	}
	if resp.Id != 2 {
		panic(fmt.Sprintf(`EXPECTED: resp.Id = 2. RECIEVED: resp.Id = %d`, resp.Id))
	}
	if resp.Nickname != "pumpalumpa" {
		panic(fmt.Sprintf(`EXPECTED: resp.Nickname = "pumpalumpa". RECIEVED: resp.Nickname = %s`, resp.Nickname))
	}
	if resp.Password != "1029384756" {
		panic(fmt.Sprintf(`EXPECTED: resp.Password = "1029384756. RECIEVED: resp.Password = %s"`, resp.Password))
	}
	if !Con.ckeckNewEmployee(resp) {
		panic("The new employee wasn't added to the database")
	}
	log.Print("Test testSignIn has been finished")
}

func testDelete() {
	defer Con.resetSequence()
	defer Con.deleteEmployees()

	Con.createEmployee()

	log.Print("Test testDelete has been started")
	del := &apptype.Delete{
		Id:       1,
		Nickname: "mama-miya",
		Passwrod: "12345678",
	}
	body, err := json.Marshal(del)
	if err != nil {
		panic(err)
	}
	answer := deleteAuthDelete(body)
	if answer != "Employee has been deleted" {
		panic(fmt.Sprintf(`EXPECTED: "Employee has been deleted". RECIEVED: %s`, answer))
	}
	if !Con.ckeckDeleteEmployee(del.Id) {
		panic("The employee wasn't deleted")
	}
	log.Print("Test testDelete has been finished")
}

func changeNickname() {
	defer Con.resetSequence()
	defer Con.deleteEmployees()

	Con.createEmployee()

	log.Print("Test changeNickname has started")
	change := &apptype.Change{
		Id:          1,
		Nickname:    "mama-miya",
		Password:    "12345678",
		WhichChange: "nickname",
		ChangeData:  "l1qwie",
	}
	body, err := json.Marshal(change)
	if err != nil {
		panic(err)
	}
	answer := putAuthChange(body)
	if answer != "The employee's data has been just successfuly changed" {
		panic(fmt.Sprintf(`EXPECTED: answer = "The employee's data has been just successfuly changed". RECIEVED: answer = %s`, answer))
	}
	if !Con.checkChangedNicknameEmployee("l1qwie") {
		panic("The employee's nickname wasn't changed")
	}
	log.Print("Test changeNickname has finished")
}

func changePassword() {
	defer Con.resetSequence()
	defer Con.deleteEmployees()

	Con.createEmployee()

	log.Print("Test changePassword has started")
	change := &apptype.Change{
		Id:          1,
		Nickname:    "mama-miya",
		Password:    "12345678",
		WhichChange: "password",
		ChangeData:  "1111",
	}
	body, err := json.Marshal(change)
	if err != nil {
		panic(err)
	}
	answer := putAuthChange(body)
	if answer != "The employee's data has been just successfuly changed" {
		panic(fmt.Sprintf(`EXPECTED: answer = "The employee's data has been just successfuly changed". RECIEVED: answer = %s`, answer))
	}
	if !Con.checkChangedPasswordEmployee("1111") {
		panic("The employee's pasword wasn't changed")
	}
	log.Print("Test changePassword has finished")
}

func testChange() {
	changeNickname()

	changePassword()
}

func StartAuthTest() {
	var err error
	log.Print("Test StartAuthTest has beed started")
	Con = new(TestConnection)
	Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	testLogIn()
	testSingIn()
	testDelete()
	testChange()
	log.Print("Test StartAuthTest has been finished")
}
