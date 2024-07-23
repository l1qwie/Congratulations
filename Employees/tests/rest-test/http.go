package resttest

import (
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

	"Employees/api/rest"
	"Employees/apptype"
	"Employees/tests/redis"
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

// Обращение к endPoint'у post-запрос
func postEmployees(body []byte, whatdo string, id int) string {
	var res string
	client := createClient()

	encryptedBody, err := encryptData(body, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt data: %s", err))
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://localhost:8099/congratulations/employees/%s/%d", whatdo, id), bytes.NewBuffer(encryptedBody))
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

// Обращение к endPoint'у get-запрос
func getEmployees(id, limit int) []*apptype.Employee {
	var employees []*apptype.Employee

	resp, err := http.Get(fmt.Sprintf("https://localhost:8099/congratulations/employees/%d/%d", id, limit))
	if err != nil {
		panic(fmt.Sprintf("Failed to make HTTP request: %s", err))
	}
	defer resp.Body.Close()

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}
	log.Printf("Response (client part): %v", respbody)

	decryptedMessage, err := decryptData(respbody, apptype.SymKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to decrypt response message: %s", err))
	}

	err = json.Unmarshal(decryptedMessage, &employees)
	if err != nil {
		var errstruct rest.Err
		err = json.Unmarshal(decryptedMessage, &errstruct)
		if err != nil {
			var resstr string
			err = json.Unmarshal(decryptedMessage, &resstr)
			if err != nil {
				panic(fmt.Sprintf("Failed to unmarshal response: %s", err))
			}
			if resstr != "" {
				panic(fmt.Sprintf("Received a string from the response: %s", resstr))
			}
		}
		if errstruct.Error != "" {
			panic(fmt.Sprintf("Error: %s", errstruct.Error))
		}
	}

	log.Print("Response:", string(decryptedMessage))
	return employees
}

// panic(fmt.Sprintf(`Expected: . Recieved: `))

// Проверка данных, которые пришли из get-запроса
func checkAnswersGetEmployees(employees []*apptype.Employee) {
	if employees[0].Id != 111 {
		panic(fmt.Sprintf("Expected: employees[0].Id = 111. Recieved: employees[0].Id = %d", employees[0].Id))
	}
	if employees[1].Id != 122 {
		panic(fmt.Sprintf("Expected: employees[1].Id = 122. Recieved: employees[1].Id = %d", employees[1].Id))
	}
	if employees[2].Id != 172 {
		panic(fmt.Sprintf("Expected: employees[2].Id = 172. Recieved: employees[2].Id = %d", employees[2].Id))
	}
	if employees[0].Name != "John" {
		panic(fmt.Sprintf(`Expected: employees[0].Name = "John". Recieved: employees[0].Name = "%s"`, employees[0].Name))
	}
	if employees[1].Name != "Moisha" {
		panic(fmt.Sprintf(`Expected: employees[1].Name = "Moisha". Recieved: employees[1].Name = "%s"`, employees[1].Name))
	}
	if employees[2].Name != "Arkadii" {
		panic(fmt.Sprintf(`Expected: employees[2].Name = "Arkadii". Recieved: employees[2].Name = "%s"`, employees[2].Name))
	}
	if employees[0].Nickname != "John1234" {
		panic(fmt.Sprintf(`Expected: employees[0].Nickname = "John1234". Recieved: employees[0].Nickname = "%s"`, employees[0].Nickname))
	}
	if employees[1].Nickname != "l1qwie" {
		panic(fmt.Sprintf(`Expected: employees[1].Nickname = "l1qwie". Recieved: employees[1].Nickname = "%s"`, employees[1].Nickname))
	}
	if employees[2].Nickname != "nagibator_2000" {
		panic(fmt.Sprintf(`Expected: employees[2].Nickname = "nagibator_2000". Recieved: employees[2].Nickname = "%s"`, employees[2].Nickname))
	}
	if employees[0].Email != "example@ya.ru" {
		panic(fmt.Sprintf(`Expected: employees[0].Email = "example@ya.ru". Recieved: employees[0].Email = "%s"`, employees[0].Email))
	}
	if employees[1].Email != "example@gmail.com" {
		panic(fmt.Sprintf(`Expected: employees[1].Email = "example@gmail.com". Recieved: employees[1].Email = "%s"`, employees[1].Email))
	}
	if employees[2].Email != "example@yahoo.com" {
		panic(fmt.Sprintf(`Expected: employees[2].Email = "example@yahoo.com". Recieved: employees[2].Email = "%s"`, employees[2].Email))
	}
	if employees[0].Birthday != "1990-03-15" {
		panic(fmt.Sprintf(`Expected: employees[0].Birthday = "1990-03-15". Recieved: employees[0].Birthday = "%s"`, employees[0].Birthday))
	}
	if employees[1].Birthday != "2003-02-20" {
		panic(fmt.Sprintf(`Expected: employees[1].Birthday = "2003-02-20". Recieved: employees[1].Birthday = "%s"`, employees[1].Birthday))
	}
	if employees[2].Birthday != "2016-09-15" {
		panic(fmt.Sprintf(`Expected: employees[2].Birthday = "2016-09-15". Recieved: employees[2].Birthday = "%s"`, employees[2].Birthday))
	}
}

// Тест для получения списка из работников с дефолтным значением 15 (максимальное количество работников за 1 запрос)
func testJustPullEmployees() {
	log.Print("testJustPullEmployees has just started")
	employees := getEmployees(0, 0)
	checkAnswersGetEmployees(employees)
	log.Print("testJustPullEmployees has successfuly finished")
}

// Тест для получения списка из трех работников
func testPullThreeEmployees() {
	log.Print("testPullThreeEmployees has just started")
	employees := getEmployees(0, 3)
	checkAnswersGetEmployees(employees)
	log.Print("testPullThreeEmployees has successfuly finished")
}

// Тест для получения списка работников с дефолтным значение 15 на которых подписан 1 (всего 4)
func testPullSubEmpToEmployees() {
	log.Print("testPullSubEmpToEmployees has just started")
	employees := getEmployees(111, 0)
	checkAnswersGetEmployees(employees)
	log.Print("testPullSubEmpToEmployees has successfuly finished")
}

// Тест для получения списка из 3 работников, на которых подписан 1 (всего 4)
func testPullThreeSubEmpToEmployees() {
	log.Print("testPullThreeSubEmpToEmployees has just started")
	employees := getEmployees(111, 3)
	checkAnswersGetEmployees(employees)
	log.Print("testPullThreeSubEmpToEmployees has successfuly finished")
}

// Тестирует endPoint c get-реквестом
func testGetEmployees() {
	log.Print("testGetEmployees has just started")
	testJustPullEmployees()
	testPullThreeEmployees()
	testPullSubEmpToEmployees()
	testPullThreeSubEmpToEmployees()
	log.Print("testGetEmployees has successfuly finished")
}

// Тест передачи данных для удаления работника из бд
func testDeleteEmployee(TRCL *redis.TestRedClient) {
	req := new(apptype.Employee)
	req.Id = 111
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	answer := postEmployees(body, "delete", 0)
	if answer != "The employee has been updated" {
		panic(fmt.Sprintf(`Expected: answer = "The employee has been updated". Recieved: answer = "%s"`, answer))
	}
	if TRCL.CheckDeletedEmployee("111") {
		panic("The employee wasn't deleted")
	}
}

// Тест передачи данных изменения данных уже существуюшего работника (удаление прошлых данных и создание новых)
func testChangeEmployee(TRCL *redis.TestRedClient) {
	req := &apptype.Employee{
		Id:       1145,
		Name:     "Изя",
		Nickname: "easy",
		Email:    "ya.ru@ya.ru",
		Birthday: "12-04-1978",
	}
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	answer := postEmployees(body, "update", 111)
	if answer != "The employee has been updated" {
		panic(fmt.Sprintf(`Expected: answer = "The employee has been updated". Recieved: answer = "%s"`, answer))
	}
	if !TRCL.CheckUpdatedOrNewEmployee(req, true, "111") {
		panic("The employee wasn't updated")
	}
}

// Тест передачи новых данных для нового работника
func testNewEmployee(TRCL *redis.TestRedClient) {
	req := &apptype.Employee{
		Id:       1145,
		Name:     "Изя",
		Nickname: "easy",
		Email:    "ya.ru@ya.ru",
		Birthday: "12-04-1978",
	}
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	answer := postEmployees(body, "new", 0)
	if answer != "The employee has been updated" {
		panic(fmt.Sprintf(`Expected: answer = "The employee has been updated". Recieved: answer = "%s"`, answer))
	}
	if !TRCL.CheckUpdatedOrNewEmployee(req, false, "") {
		panic("The new employee wasn't added")
	}
}

func testSubEmployeeToEmployee(TRCL *redis.TestRedClient) {
	req := new(apptype.Employee)
	req.Id = 111
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	answer := postEmployees(body, "sub", 199)
	if answer != "The employee has been updated" {
		panic(fmt.Sprintf(`Expected: answer = "The employee has been updated". Recieved: answer = "%s"`, answer))
	}
	if !TRCL.CheckSubToEmployee(req.Id, 199) {
		panic("The employee wasn't subed to another employee")
	}
}

func testUnSubEmployeeFromEmployee(TRCL *redis.TestRedClient) {
	req := new(apptype.Employee)
	req.Id = 111
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	answer := postEmployees(body, "unsub", 122)
	if answer != "The employee has been updated" {
		panic(fmt.Sprintf(`Expected: answer = "The employee has been updated". Recieved: answer = "%s"`, answer))
	}
	if TRCL.CheckSubToEmployee(req.Id, 122) {
		panic("The employee wasn't unsubed to another employee")
	}
}

// Тестирует endPoint c post-реквестом
func testPostEmployees(TRCL *redis.TestRedClient) {
	log.Print("testPostEmployees has just started")
	testDeleteEmployee(TRCL)
	testChangeEmployee(TRCL)
	testNewEmployee(TRCL)
	testSubEmployeeToEmployee(TRCL)
	testUnSubEmployeeFromEmployee(TRCL)
	log.Print("testPostEmployees has successfuly finished")
}

// Начинает тесты для микросервиса Employees
func StartEmployeeTests() {
	var err error
	log.Print("Tests StartEmployeeTests() started")
	TRCL := new(redis.TestRedClient)
	TRCL.Cl, err = redis.AddClient()
	if err != nil {
		panic(err)
	}
	TRCL.InitEmployees()
	TRCL.SubscribeEmployees()
	defer TRCL.DeleteEmployees()

	testGetEmployees()
	testPostEmployees(TRCL)
	log.Print("Tests StartEmployeeTests() finished")
}
