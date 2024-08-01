package kafkatest

import (
	"Subscribe/apptype"
	resttest "Subscribe/tests/rest-test"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	maxRetries   int           = 5
	retryBackoff time.Duration = 2 * time.Second
	topic        string        = "employee-sub"
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

func producer(subscriber, subtoid int, whatdo string) {
	log.Printf("Got into producer() with params subscriber: %d, subtoid: %d, whatdo: %s", subscriber, subtoid, whatdo)
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
		Id:       subscriber,
		SecondId: subtoid,
		WhatDo:   whatdo,
	}
	jb, err := json.Marshal(kafkaemployee)
	if err != nil {
		log.Printf("KAFKA ERROR TellChanges(): %s", err)
		return
	}
	send(jb, producer, 0)
	log.Print("Got out of producer()")
}

func testSub() {
	log.Print("Test testSub() has started")
	defer resttest.Con.ResetSequence()
	defer resttest.Con.DeleteEmployees()
	defer resttest.Con.DeleteSubscribers()

	resttest.Con.CreateEmployee()

	producer(1, 2, "sub")
	time.Sleep(time.Second * 2)
	if !resttest.Con.CheckSubedEmployees(1, 2) {
		panic("The employee wasn't subscribed to another")
	}
	log.Print("Test testSub() has finished")
}

func testUnsub() {
	log.Print("Test testUnsub() has started")
	defer resttest.Con.ResetSequence()
	defer resttest.Con.DeleteEmployees()
	defer resttest.Con.DeleteSubscribers()

	resttest.Con.CreateEmployee()
	resttest.Con.SubscribeEmployeeToEmployee()

	producer(1, 2, "unsub")
	time.Sleep(time.Second * 2)
	if resttest.Con.CheckSubedEmployees(1, 2) {
		panic("The employee wasn't unsubscribed to another")
	}
	log.Print("Test testUnsub() has finished")
}

func StartSubKafkaTest() {
	time.Sleep(time.Second * 5)
	testSub()
	testUnsub()
}
