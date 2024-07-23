package consumer

import (
	"Employees/app"
	"Employees/apptype"
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
)

func whichWay(kafkaemployee *apptype.KafkaEmployee) {
	var err error
	employee := new(apptype.Employee)
	employee.Id = kafkaemployee.Id
	employee.Name = kafkaemployee.Name
	employee.Nickname = kafkaemployee.Nickname
	employee.Email = kafkaemployee.Email
	employee.Birthday = kafkaemployee.Birthday
	if kafkaemployee.WhatDo == add {
		err = app.Client.NewEmpl(employee)
	} else if kafkaemployee.WhatDo == update {
		err = app.Client.UpdEmpl(employee, kafkaemployee.SecondId)
	} else if kafkaemployee.WhatDo == delete {
		err = app.Client.DeleleEmpl(employee.Id)
	} else if kafkaemployee.WhatDo == sub {
		err = app.Client.AddSub(employee.Id, kafkaemployee.SecondId)
	} else if kafkaemployee.WhatDo == unsub {
		err = app.Client.UnSub(employee.Id, kafkaemployee.SecondId)
	}
	if err != nil {
		log.Printf("!!!ERROR!!! : %s", err)
	}
}

func reader(top string, partcons sarama.PartitionConsumer) {
	kafkaemployee := new(apptype.KafkaEmployee)
	for {
		select {
		case msg := <-partcons.Messages():
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
	brokers := []string{"kafka:9092"}
	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer: %s", err))
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, Employee, sarama.OffsetNewest)
	if err != nil {
		panic(fmt.Sprintf("Failed to start consumer partition add: %s", err))
	}
	reader(topic, partitionConsumer)
}
