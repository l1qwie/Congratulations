package tests

import (
	"Notifications/app"
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
	"reflect"
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
		Birthday:      "1980-02-20",
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

func testFedya(memebers map[*apptype.Employee][]*apptype.Employee) {
	fedya := &apptype.Employee{
		Id:       131,
		Name:     "Fedor",
		Nickname: "fedor999",
		Birthday: "1919-09-09",
	}
	misha := &apptype.Employee{
		Id:       127,
		Name:     "Misha",
		Nickname: "misha228",
		Email:    "trashfly@ya.ru",
	}
	// katya := &apptype.Employee{
	// 	Id:       128,
	// 	Name:     "Katia",
	// 	Nickname: "krasotka_katusha",
	// 	Email:    "example-example@gmail.com",
	// }
	// alesha := &apptype.Employee{
	// 	Id:       129,
	// 	Name:     "Alesha",
	// 	Nickname: "alexsei",
	// 	Email:    "me@yahoo.com",
	// }
	fedyaval, fedyaok := memebers[fedya]
	fedyacompare := []*apptype.Employee{misha} //katya, alesha
	if fedyaok {
		log.Print("fedyaok is true")
		for _, val := range fedyaval {
			for i, val2 := range fedyacompare {
				if !reflect.DeepEqual(val, val2) {
					if i+1 == len(fedyacompare) {
						panic(fmt.Sprintf("Data in fedyaval: %v couldn't be found in fedyacompare: %v", val, fedyacompare))
					}
				}
				log.Printf("Data in fedyaval: %v. Data in fedyacompare: %v", val, val2)
			}
		}
	}
}

func testOleja(memebers map[*apptype.Employee][]*apptype.Employee) {
	oleja := &apptype.Employee{
		Id:       126,
		Name:     "Oleja",
		Nickname: "oleja_krut",
		Birthday: "1980-09-20",
	}
	misha := &apptype.Employee{
		Id:       127,
		Name:     "Misha",
		Nickname: "misha228",
		Email:    "trashfly@ya.ru",
	}
	nastya := &apptype.Employee{
		Id:       130,
		Name:     "Nastya",
		Nickname: "princes",
		Email:    "letmefly@gmail.com",
	}
	olejaval, olejaok := memebers[oleja]
	olejacompare := []*apptype.Employee{misha, nastya}
	if olejaok {
		log.Print("olejaok is true")
		for _, val := range olejaval {
			for i, val2 := range olejacompare {
				if !reflect.DeepEqual(val, val2) {
					if i+1 == len(olejacompare) {
						panic(fmt.Sprintf("Data in fedyaval: %v couldn't be found in fedyacompare: %v", val, olejacompare))
					}
				}
				log.Printf("Data in fedyaval: %v. Data in fedyacompare: %v", val, val2)
			}
		}
	}
}

func testMatvei(memebers map[*apptype.Employee][]*apptype.Employee) {
	matvei := &apptype.Employee{
		Id:       132,
		Name:     "Matvei",
		Nickname: "cool_juice",
		Birthday: "1980-10-03",
	}
	nastya := &apptype.Employee{
		Id:       130,
		Name:     "Nastya",
		Nickname: "princes",
		Email:    "letmefly@gmail.com",
	}
	misha := &apptype.Employee{
		Id:       127,
		Name:     "Misha",
		Nickname: "misha228",
		Email:    "trashfly@ya.ru",
	}
	// katya := &apptype.Employee{
	// 	Id:       128,
	// 	Name:     "Katia",
	// 	Nickname: "krasotka_katusha",
	// 	Email:    "example-example@gmail.com",
	// }
	// alesha := &apptype.Employee{
	// 	Id:       129,
	// 	Name:     "Alesha",
	// 	Nickname: "alexsei",
	// 	Email:    "me@yahoo.com",
	// }
	matveival, matveiok := memebers[matvei]
	matveicompare := []*apptype.Employee{nastya, misha} //katya, alesha
	if matveiok {
		log.Print("matveiok is true")
		for _, val := range matveival {
			for i, val2 := range matveicompare {
				if !reflect.DeepEqual(val, val2) {
					if i+1 == len(matveicompare) {
						panic(fmt.Sprintf("Data in fedyaval: %v couldn't be found in fedyacompare: %v", val, matveicompare))
					}
				}
				log.Printf("Data in fedyaval: %v. Data in fedyacompare: %v", val, val2)
			}
		}
	}
}

func testBogdan(memebers map[*apptype.Employee][]*apptype.Employee) {
	bogdan := &apptype.Employee{
		Id:       125,
		Name:     "Bogdan",
		Nickname: "l1qwie",
		Birthday: "2003-02-20",
	}
	nastya := &apptype.Employee{
		Id:       130,
		Name:     "Nastya",
		Nickname: "princes",
		Email:    "letmefly@gmail.com",
	}
	// katya := &apptype.Employee{
	// 	Id:       128,
	// 	Name:     "Katia",
	// 	Nickname: "krasotka_katusha",
	// 	Email:    "example-example@gmail.com",
	// }
	bogdanval, bogdanok := memebers[bogdan]
	bogdancompare := []*apptype.Employee{nastya} //katya
	if bogdanok {
		log.Print("matveiok is true")
		for _, val := range bogdanval {
			for i, val2 := range bogdancompare {
				if !reflect.DeepEqual(val, val2) {
					if i+1 == len(bogdancompare) {
						panic(fmt.Sprintf("Data in fedyaval: %v couldn't be found in fedyacompare: %v", val, bogdancompare))
					}
				}
				log.Printf("Data in fedyaval: %v. Data in fedyacompare: %v", val, val2)
			}
		}
	}
}

func UnitTest() {
	var err error
	app.Con = new(app.Connection)
	app.Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	Con = new(TestConnection)
	Con.DB, err = apptype.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	defer Con.DeleteEmployees()
	defer Con.DeleteSubscriptions()

	Con.Create8Employees()
	Con.SubEmployeesToEmployees()

	memebers := app.Con.FindWhoShouldBeNotified()
	log.Print(len(memebers))
	for key, val := range memebers {
		log.Printf("Key: %v", key)
		for _, val2 := range val {
			log.Printf("Val: %v", val2)
		}
	}
	testFedya(memebers)
	testOleja(memebers)
	testMatvei(memebers)
	testBogdan(memebers)
}
