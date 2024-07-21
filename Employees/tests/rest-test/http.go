package resttest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/l1qwie/Congratulations/Employees/apptype"
	"github.com/l1qwie/Congratulations/Employees/tests/redis"
)

var key []byte

// Creates HTTPS client with settings
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

func getEmployees(id, limit int) []*apptype.Employee {
	resp, err := http.Get(fmt.Sprintf("https://localhost:8099/congratulations/employees/%d/%d", id, limit))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response: %s", err))
	}
	log.Printf("Response: %s", string(respbody))
	employees := make([]*apptype.Employee, limit)
	err = json.Unmarshal(respbody, &employees)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal response: %s", err))
	}
	return employees
}

// panic(fmt.Sprintf(`Expected: . Recieved: `))
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

func testJustPullEmployees() {
	log.Print("testJustPullEmployees has just started")
	employees := getEmployees(0, 0)
	checkAnswersGetEmployees(employees)
	log.Print("testJustPullEmployees has successfuly finished")
}

func testPullThreeEmployees() {
	log.Print("testPullThreeEmployees has just started")
	employees := getEmployees(0, 3)
	checkAnswersGetEmployees(employees)
	log.Print("testPullThreeEmployees has successfuly finished")
}

func testPullSubEmpToEmployees() {
	log.Print("testPullSubEmpToEmployees has just started")
	employees := getEmployees(111, 0)
	checkAnswersGetEmployees(employees)
	log.Print("testPullSubEmpToEmployees has successfuly finished")
}

func testPullThreeSubEmpToEmployees() {
	log.Print("testPullThreeSubEmpToEmployees has just started")
	employees := getEmployees(111, 3)
	checkAnswersGetEmployees(employees)
	log.Print("testPullThreeSubEmpToEmployees has successfuly finished")
}

func testGetEmployees() {
	log.Print("testGetEmployees has just started")
	testJustPullEmployees()
	testPullThreeEmployees()
	testPullSubEmpToEmployees()
	testPullThreeSubEmpToEmployees()
	log.Print("testGetEmployees has successfuly finished")
}

func StartEmployeeTests() {
	var err error
	TRCL := new(redis.TestRedClient)
	TRCL.Cl, err = redis.AddClient()
	if err != nil {
		panic(err)
	}
	TRCL.InitEmployees()
	TRCL.SubscribeEmployees()
	defer TRCL.DeleteEmployees()

	testGetEmployees()
}
