package kafka

import (
	//"Employees/api/kafka/producer"
	"Employees/apptype"
	"Employees/tests/redis"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	maxRetries   int           = 5
	retryBackoff time.Duration = 2 * time.Second
	TestTopic    string        = "employee-redis"
)

func send(jd []byte, producer sarama.AsyncProducer, partition int32) {
	msg := &sarama.ProducerMessage{
		Topic:     TestTopic,
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

func producer(employee *apptype.Employee, whatdo string, secondid int) {
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

	kafkaemployee := &apptype.KafkaEmployee{
		Id:       employee.Id,
		Name:     employee.Name,
		Nickname: employee.Nickname,
		Email:    employee.Email,
		Birthday: employee.Birthday,
		WhatDo:   whatdo,
		SecondId: secondid,
	}

	jb, err := json.Marshal(kafkaemployee)
	if err != nil {
		log.Printf("KAFKA ERROR TellChanges(): %s", err)
		return
	}

	send(jb, producer, 0)
}

func testNew(TRCL *redis.TestRedClient) {
	log.Print("testNew has been started")
	employee := &apptype.Employee{
		Id:       168,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer(employee, "new", 0)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckUpdatedOrNewEmployee(employee, false, "") {
		panic("The new employee wasn't added")
	}
	log.Print("testNew has been finished")
}

func testUpdate(TRCL *redis.TestRedClient) {
	log.Print("testUpdate has been started")
	employee := &apptype.Employee{
		Id:       168,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer(employee, "update", 111)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckUpdatedOrNewEmployee(employee, true, "111") {
		panic("The employee wasn't updated")
	}
	log.Print("testUpdate has been finished")
}

func testDelete(TRCL *redis.TestRedClient) {
	log.Print("testDelete has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer(employee, "delete", 0)
	time.Sleep(time.Second / 10)
	if TRCL.CheckDeletedEmployee("111") {
		panic("The employee wasn't deleted")
	}
	log.Print("testDelete has been finished")
}

func testSub(TRCL *redis.TestRedClient) {
	log.Print("testSub has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer(employee, "sub", 199)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckSubToEmployee(employee.Id, 199) {
		panic("The employee wasn't subed to another employee")
	}
	log.Print("testSub has been finished")
}

func testUnsub(TRCL *redis.TestRedClient) {
	log.Print("testUnsub has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer(employee, "unsub", 111)
	time.Sleep(time.Second / 10)
	if TRCL.CheckSubToEmployee(employee.Id, 111) {
		panic("The employee wasn't unsubed to another employee")
	}
	log.Print("testUnsub has been finished")
}

func TestConsumer() {
	log.Print("TestConsumer has been started")
	var err error
	TRCL := new(redis.TestRedClient)
	TRCL.Cl, err = redis.AddClient()
	if err != nil {
		panic(err)
	}

	TRCL.InitEmployees()
	TRCL.SubscribeEmployees()
	defer TRCL.DeleteEmployees()

	time.Sleep(time.Second / 10)

	testNew(TRCL)
	testUpdate(TRCL)
	testDelete(TRCL)
	testSub(TRCL)
	testUnsub(TRCL)
	log.Print("TestConsumer has been finished")
}
