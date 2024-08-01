package resttest

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"Subscribe/apptype"
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

// Расшифровывает данные от сервера
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

func putSubs(method string, subscriber, subto int) string {
	client := createClient()
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://localhost:8019/congratulations/subscriptions/%s/%d/%d", method, subscriber, subto), bytes.NewBuffer(nil))
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

func testSub() {
	defer Con.DeleteEmployees()
	defer Con.DeleteSubscribers()
	defer Con.ResetSequence()

	Con.CreateEmployee()

	answer := putSubs("sub", 1, 2)
	if answer != "You've just subscribed to 2" {
		panic(fmt.Sprintf(`EXPECTED: "You've just subscribed to 2". RECIEVED: %s`, answer))
	}
	if !Con.CheckSubedEmployees(1, 2) {
		panic("The employee wasn't subscribed to another employee")
	}
}

func testUnsub() {
	defer Con.DeleteEmployees()
	defer Con.DeleteSubscribers()
	defer Con.ResetSequence()

	Con.CreateEmployee()
	Con.SubscribeEmployeeToEmployee()

	answer := putSubs("unsub", 1, 2)
	if answer != "You've just unsubscribed from 2" {
		panic(fmt.Sprintf(`EXPECTED: "You've just unsubscribed from 2". RECIEVED: %s`, answer))
	}
	if Con.CheckSubedEmployees(1, 2) {
		panic("The employee wasn't unsubscribed to another employee")
	}
}

func StartSubTest() {
	testSub()
	testUnsub()
}
