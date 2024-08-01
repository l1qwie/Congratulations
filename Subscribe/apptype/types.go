package apptype

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	docHost     = "congratulations-postgres"
	docPort     = "5432"
	docUsername = "postgres"
	docPass     = "postgres"
	docDbname   = "postgres"
	docSslmode  = "disable"
)

var SymKey []byte

type KafkaEmployee struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	WhatDo   string `json:"whatdo"`
	SecondId int    `json:"secondid"`
}

func docConnect() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		docHost,
		docPort,
		docUsername,
		docPass,
		docDbname,
		docSslmode)
}

func ConnectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", docConnect())
	if err != nil {
		log.Print(err)
	} else {
		err = db.Ping()
		if err != nil {
			log.Print(err)
		}
	}
	return db, err
}
