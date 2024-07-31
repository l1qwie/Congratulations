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

func testFedya(memebers map[apptype.Employee][]*apptype.Employee) {
	var found bool
	fedya := apptype.Employee{
		Id:       131,
		Name:     "Fedor",
		Nickname: "fedor999",
		Birthday: "1919-08-09",
	}
	misha := &apptype.Employee{
		Id:       127,
		Name:     "Misha",
		Nickname: "misha228",
		Email:    "trashfly@ya.ru",
	}
	katya := &apptype.Employee{
		Id:       128,
		Name:     "Katia",
		Nickname: "krasotka_katusha",
		Email:    "example-example@gmail.com",
	}
	alesha := &apptype.Employee{
		Id:       129,
		Name:     "Alesha",
		Nickname: "alexsei",
		Email:    "me@yahoo.com",
	}
	log.Printf("Fedya: %v", fedya)
	fedyaval, fedyaok := memebers[fedya]
	fedyacompare := []*apptype.Employee{misha, katya, alesha}
	if fedyaok {
		log.Print("fedyaok is true")
		if len(fedyaval) == len(fedyacompare) {
			for _, val := range fedyaval {
				for i := 0; i < len(fedyacompare) && !found; i++ {
					if !reflect.DeepEqual(*val, *fedyacompare[i]) {
						if i+1 == len(fedyacompare) {
							if !found {
								panic(fmt.Sprintf("Data in fedyaval: %v couldn't be found in fedyacompare: %v", *val, fedyacompare))
							}
						}
					} else {
						found = true
					}
					log.Printf("Data in fedyaval: %v. Data in fedyacompare: %v. found = %v", *val, *fedyacompare[i], found)
				}
				found = false
			}
		} else {
			panic(fmt.Sprintf("len(fedyaval) = EXPECTED: 3, RECIEVED: %d\nlen(fedyacompare) = EXPECTED: 3, RECIEVED: %d", len(fedyaval), len(fedyacompare)))
		}
	} else {
		panic("fedyaok = false")
	}
}

func testOleja(memebers map[apptype.Employee][]*apptype.Employee) {
	var found bool
	oleja := apptype.Employee{
		Id:       126,
		Name:     "Oleja",
		Nickname: "oleja_krut",
		Birthday: "1980-08-20",
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
	olejasubs, olejaok := memebers[oleja]
	olejacompare := []*apptype.Employee{misha, nastya}
	if olejaok {
		log.Print("fedyaok is true")
		if len(olejasubs) == len(olejacompare) {
			for _, val := range olejasubs {
				for i := 0; i < len(olejacompare) && !found; i++ {
					if !reflect.DeepEqual(*val, *olejacompare[i]) {
						if i+1 == len(olejacompare) {
							if !found {
								panic(fmt.Sprintf("Data in olejasubs: %v couldn't be found in olejacompare: %v", *val, olejacompare))
							}
						}
					} else {
						found = true
					}
					log.Printf("Data in olejasubs: %v. Data in olejacompare: %v. found = %v", *val, *olejacompare[i], found)
				}
				found = false
			}
		} else {
			panic(fmt.Sprintf("len(olejasubs) = EXPECTED: 2, RECIEVED: %d\nlen(olejacompare) = EXPECTED: 2, RECIEVED: %d", len(olejasubs), len(olejacompare)))
		}
	} else {
		panic("olejaok = false")
	}
}

func testMatvei(memebers map[apptype.Employee][]*apptype.Employee) {
	var found bool
	matvei := apptype.Employee{
		Id:       132,
		Name:     "Matvei",
		Nickname: "cool_juice",
		Birthday: "1980-08-03",
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
	katya := &apptype.Employee{
		Id:       128,
		Name:     "Katia",
		Nickname: "krasotka_katusha",
		Email:    "example-example@gmail.com",
	}
	alesha := &apptype.Employee{
		Id:       129,
		Name:     "Alesha",
		Nickname: "alexsei",
		Email:    "me@yahoo.com",
	}
	matveisubs, matveiok := memebers[matvei]
	matveicompare := []*apptype.Employee{nastya, misha, katya, alesha}
	if matveiok {
		log.Print("fedyaok is true")
		if len(matveisubs) == len(matveicompare) {
			for _, val := range matveisubs {
				for i := 0; i < len(matveicompare) && !found; i++ {
					if !reflect.DeepEqual(*val, *matveicompare[i]) {
						if i+1 == len(matveicompare) {
							if !found {
								panic(fmt.Sprintf("Data in matveisubs: %v couldn't be found in matveicompare: %v", *val, matveicompare))
							}
						}
					} else {
						found = true
					}
					log.Printf("Data in matveisubs: %v. Data in matveicompare: %v. found = %v", *val, matveicompare[i], found)
				}
				found = false
			}
		} else {
			panic(fmt.Sprintf("len(matveisubs) = EXPECTED: 4, RECIEVED: %d\nlen(matveicompare) = EXPECTED: 4, RECIEVED: %d", len(matveisubs), len(matveicompare)))
		}
	} else {
		panic("matveiok = false")
	}
}

func testBogdan(memebers map[apptype.Employee][]*apptype.Employee) {
	var found bool
	bogdan := apptype.Employee{
		Id:       125,
		Name:     "Bogdan",
		Nickname: "l1qwie",
		Birthday: "2003-07-30",
	}
	nastya := &apptype.Employee{
		Id:       130,
		Name:     "Nastya",
		Nickname: "princes",
		Email:    "letmefly@gmail.com",
	}
	katya := &apptype.Employee{
		Id:       128,
		Name:     "Katia",
		Nickname: "krasotka_katusha",
		Email:    "example-example@gmail.com",
	}
	bogdansubs, bogdanok := memebers[bogdan]
	bogdancompare := []*apptype.Employee{nastya, katya}
	if bogdanok {
		log.Print("fedyaok is true")
		if len(bogdansubs) == len(bogdancompare) {
			for _, val := range bogdansubs {
				for i := 0; i < len(bogdancompare) && !found; i++ {
					if !reflect.DeepEqual(*val, *bogdancompare[i]) {
						if i+1 == len(bogdancompare) {
							if !found {
								panic(fmt.Sprintf("Data in bogdansubs: %v couldn't be found in bogdancompare: %v", *val, bogdancompare))
							}
						}
					} else {
						found = true
					}
					log.Printf("Data in bogdansubs: %v. Data in bogdancompare: %v. found = %v", *val, bogdancompare[i], found)
				}
				found = false
			}
		} else {
			panic(fmt.Sprintf("len(bogdansubs) = EXPECTED: 2, RECIEVED: %d\nlen(bogdancompare) = EXPECTED: 2, RECIEVED: %d", len(bogdansubs), len(bogdancompare)))
		}
	} else {
		panic("bogdanok = false")
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
