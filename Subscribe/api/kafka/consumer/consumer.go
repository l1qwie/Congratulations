package consumer

import (
	"Subscribe/app"
	"Subscribe/apptype"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

const (
	subb  string = "sub"
	unsub string = "unsub"
	topic string = "employee-sub"
)

func whichWay(kafkaemployee *apptype.KafkaEmployee) {
	log.Printf("Got into whichWay() with param kafkaemployee: %v", kafkaemployee)
	var (
		err                          error
		ok, employee1ok, employee2ok bool
	)
	log.Printf("kafkaemployee.WhatDo = %s", kafkaemployee.WhatDo)
	employee1ok, err = app.Con.FindEmploee(kafkaemployee.Id)
	if employee1ok {
		employee2ok, err = app.Con.FindEmploee(kafkaemployee.SecondId)
		if employee2ok {
			ok, err = app.Con.CheckSubStatus(kafkaemployee.Id, kafkaemployee.SecondId)
			if ok && kafkaemployee.WhatDo == subb && err == nil {
				log.Print("STRAT SUBSCRIBING")
				err = app.Con.SubEmployeeToEmployee(kafkaemployee.Id, kafkaemployee.SecondId)
			} else if !ok && kafkaemployee.WhatDo == unsub && err == nil {
				log.Print("START UNSUBSCRIBING")
				err = app.Con.UnsubEmployeeFromEmployee(kafkaemployee.Id, kafkaemployee.SecondId)
			}
		}
	}
	if err != nil {
		log.Printf("!!!ERROR!!! : %s", err)
	}
	log.Print("Got out of whichWay()")
}

func reader(partcons sarama.PartitionConsumer) {
	log.Printf("Started reader() with paramms top: %s, partcons: %v", topic, partcons)
	kafkaemployee := new(apptype.KafkaEmployee)
	for {
		select {
		case msg := <-partcons.Messages():
			log.Print("Got a message")
			err := json.Unmarshal(msg.Value, &kafkaemployee)
			if err != nil {
				log.Print(err)
			}
			whichWay(kafkaemployee)
		case err := <-partcons.Errors():
			log.Printf("Error while consuming (topic: %s): %s", topic, err)
		}

	}
}

func Consumer() {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	brokers := []string{"congratulations-kafka:9092"}
	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer: %s", err))
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer partition add: %s", err))
	}
	reader(partitionConsumer)
}
