package kafkatest

import (
	"Authorization/api/kafka/producer"
	"Authorization/apptype"
	resttest "Authorization/tests/rest-test"
	"log"
	"time"
)

const (
	Auth         int32         = 1
	maxRetries   int           = 5
	retryBackoff time.Duration = 2 * time.Second
)

func testNewEmployee() {
	defer resttest.Con.ResetSequence()
	defer resttest.Con.DeleteDONOTUSE()
	defer resttest.Con.DeleteEmployees()

	log.Print("Started test testNewEmployee()")
	kafkaempl := &apptype.KafkaEmployee{
		Id:       111,
		Name:     "Miya",
		Nickname: "mama-miya",
		Email:    "123455@exapmle.com",
		Birthday: "10-09-2000",
		WhatDo:   "new",
	}
	producer.TellChanges(kafkaempl)
	time.Sleep(time.Second * 3)
	auth, employee := resttest.Con.CkeckAddedOrUpdatedEmployeeKafka(kafkaempl)
	if !auth || !employee {
		if !auth {
			log.Printf("Error during checking auth: %v", auth)
		} else if !employee {
			log.Printf("Error during checking check: %v", employee)
		}
		panic("The new employee wasn't added")
	}
	log.Print("Finished test testNewEmployee()")
}

func testUpdateEmployee() {
	log.Print("Test testUpdateEmployee has been started")
	defer resttest.Con.ResetSequence()
	defer resttest.Con.DeleteDONOTUSE()
	defer resttest.Con.DeleteEmployees()

	resttest.Con.CreateEmployee()
	resttest.Con.CreateEmployeeKafka()

	kafkaempl := &apptype.KafkaEmployee{
		Id:       199,
		Name:     "Miya",
		Nickname: "shushara",
		Email:    "12@exapmle.com",
		Birthday: "07-09-2000",
		WhatDo:   "update",
		SecondId: 1,
	}
	producer.TellChanges(kafkaempl)
	// consumer.Consumer()
	time.Sleep(time.Second * 3)
	auth, employee := resttest.Con.CkeckAddedOrUpdatedEmployeeKafka(kafkaempl)
	if !auth || !employee {
		if !auth {
			log.Printf("Error during checking auth: %v", auth)
		} else if !employee {
			log.Printf("Error during checking check: %v", employee)
		}
		panic("The new employee wasn't updated")
	}
	log.Print("Test testUpdateEmployee has been finished")
}

func testDeleteEmployee() {
	log.Print("Test testDeleteEmployee has been started")
	defer resttest.Con.ResetSequence()
	defer resttest.Con.DeleteDONOTUSE()
	defer resttest.Con.DeleteEmployees()

	resttest.Con.CreateEmployee()
	resttest.Con.CreateEmployeeKafka()

	kafkaempl := &apptype.KafkaEmployee{
		Id:       199,
		Name:     "Miya",
		Nickname: "shushara",
		Email:    "12@exapmle.com",
		Birthday: "07-09-2000",
		WhatDo:   "update",
		SecondId: 1,
	}
	producer.TellChanges(kafkaempl)
	time.Sleep(time.Second / 10)
	auth, employees, subs := resttest.Con.CheckDeletedEmployeeKafka(kafkaempl)
	if !auth || !employees || !subs {
		if !auth {
			log.Printf("Error in during checking auth: %v", auth)
		} else if !employees {
			log.Printf("Error in during checking employees: %v", employees)
		} else if !subs {
			log.Printf("Error in during checking subs: %v", subs)
		}
		panic("The employee wasn't deleted")
	}
	log.Print("Test testDeleteEmployee has been finished")
}

func StartTestConsumer() {
	testNewEmployee()
	testUpdateEmployee()
	testDeleteEmployee()
}
