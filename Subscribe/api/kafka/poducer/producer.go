package producer

import (
	"Subscribe/apptype"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	maxRetries   int           = 5
	retryBackoff time.Duration = 2 * time.Second
	topic        string        = "employee-redis"
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

func TellChanges(subscriber, subtoid int, whatdo string) {
	log.Printf("Got into producer.TellChanges() with params subscriber: %d, subtoid: %d, whatdo: %s", subscriber, subtoid, whatdo)
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
	log.Print("Got out of producer.TellChanges()")
}
