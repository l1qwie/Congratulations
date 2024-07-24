package tests

import (
	"Notifications/apptype"
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

func postNotifications(body []byte) string {
	var res string
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("POST", "https://localhost:8089/congratulations/notify", bytes.NewBuffer(encryptedBody))
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
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}

	res = string(decryptedMessage)
	log.Print("Response:", res)
	return res
}

func sendNotifier() {
	notified := &apptype.Notified{
		IdSubed:       125,
		IdBirth:       126,
		NameSubed:     "Bogdan",
		NameBirth:     "Oleja",
		NicknameSubed: "l1qwie",
		NicknameBirth: "oleja_krut",
		Email:         "trashemailbyfly@gmail.com",
		Birthday:      "20-09-1980",
	}
	body, err := json.Marshal(notified)
	if err != nil {
		panic(err)
	}
	answer := postNotifications(body)
	if answer != "The employee was notified" {
		panic(fmt.Sprintf(`Expected: The employee was notified. Recieved: %s`, answer))
	}
	if !Con.checkNotifiedInDB(125, 126) {
		panic("The employee wasn't notified according to the database")
	}
}

func TestNotifications() {
	var err error
	Con = new(TestConnection)
	Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	defer Con.DeleteEmployees()
	defer Con.DeleteSubscriptions()

	Con.CreateEmployees()
	Con.SubscribeToEmployee()

	sendNotifier()
}
