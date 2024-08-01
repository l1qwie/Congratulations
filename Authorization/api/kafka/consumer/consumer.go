package consumer

import (
	"Authorization/app"
	"Authorization/apptype"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

const (
	add      string = "new"
	update   string = "update"
	delete   string = "delete"
	sub      string = "sub"
	unsub    string = "unsub"
	Employee int32  = 0
	Auth     int32  = 1
)

func whichWay(kafkaemployee *apptype.KafkaEmployee) {
	log.Printf("Got into whichWay() with param kafkaemployee: %v", kafkaemployee)
	var err error
	con := new(app.Connection)
	con.DB, err = apptype.ConnectToDatabase()
	if err == nil {
		log.Print("The app has successfuly connected to the database")
		defer con.DB.Close()
		log.Printf("kafkaemployee.WhatDo = %s", kafkaemployee.WhatDo)
		if kafkaemployee.WhatDo == add {
			err = con.AddEmployee(kafkaemployee)
		} else if kafkaemployee.WhatDo == update {
			err = con.UpdateEmployee(kafkaemployee)
		} else if kafkaemployee.WhatDo == delete {
			err = con.DeleteEmployee(kafkaemployee.Id)
		}
	}
	if err != nil {
		log.Printf("!!!ERROR!!! : %s", err)
	}
	log.Print("Got out of whichWay()")
}

func reader(top string, partcons sarama.PartitionConsumer) {
	log.Printf("Started reader() with paramms top: %s, partcons: %v", top, partcons)
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
			log.Printf("Error while consuming (topic: %s): %s", top, err)
		}
	}
}

func Consumer() {
	topic := "employee"
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	brokers := []string{"congratulations-kafka:9092"}
	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer: %s", err))
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, Auth, sarama.OffsetNewest)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer partition add: %s", err))
	}
	reader(topic, partitionConsumer)
}
