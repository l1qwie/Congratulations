package apptype

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	SymKey []byte
	TestIP string
)

const (
	docHost     = "congratulations-postgres"
	docPort     = "5432"
	docUsername = "postgres"
	docPass     = "postgres"
	docDbname   = "postgres"
	docSslmode  = "disable"
)

type Auth struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type SignIn struct {
	Answer   string `json:"answer"`
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type Delete struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Passwrod string `json:"password"`
}

type Change struct {
	Id          int    `json:"id"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	WhichChange string `json:"which_change"`
	ChangeData  string `json:"change_data"`
}

type KafkaEmployee struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	WhatDo   string `json:"whatdo"`
	SecondId int    `json:"secondid"`
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

func docConnect() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		docHost,
		docPort,
		docUsername,
		docPass,
		docDbname,
		docSslmode)
}
