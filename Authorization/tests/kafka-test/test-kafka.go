package kafkatest

import (
	"Authorization/apptype"
	resttest "Authorization/tests/rest-test"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	maxRetries   int           = 5
	retryBackoff time.Duration = 2 * time.Second
	topic        string        = "employee-other"
)

func send(jd []byte, producer sarama.AsyncProducer, partition int32) {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(jd),
		Partition: partition,
	}
	continuworking := true
	attempt := 0
	for continuworking {
		producer.Input() <- msg

		select {
		case <-producer.Successes():
			log.Print("Message has been successfully sent")
			continuworking = false
		case err := <-producer.Errors():
			log.Printf("Failed to send message: %s", err)
			attempt++
			if attempt >= maxRetries {
				log.Printf("Max retries reached. Giving up on sending message.")
				continuworking = false
			}
			log.Printf("Retrying in %v...", retryBackoff)
			time.Sleep(retryBackoff)
		}
	}
}

func producer(kafkaemployee *apptype.KafkaEmployee) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Retry.Backoff = 100 * time.Millisecond
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Timeout = 10 * time.Second
	brokers := []string{"congratulations-kafka:9092"}

	producer, err := sarama.NewAsyncProducer(brokers, producerConfig)
	if err != nil {
		log.Printf("Failed to start producer: %s", err)
		return
	}
	defer producer.Close()

	jb, err := json.Marshal(kafkaemployee)
	if err != nil {
		log.Printf("KAFKA ERROR TellChanges(): %s", err)
		return
	}

	send(jb, producer, 0)
}

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
	producer(kafkaempl)
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
	producer(kafkaempl)
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
		WhatDo:   "delete",
		SecondId: 1,
	}
	producer(kafkaempl)
	time.Sleep(time.Second / 10)
	auth, employees, subs := resttest.Con.CheckDeletedEmployeeKafka(kafkaempl)
	if auth || employees || !subs {
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
	time.Sleep(time.Second * 3)
	testNewEmployee()
	testUpdateEmployee()
	testDeleteEmployee()
}
