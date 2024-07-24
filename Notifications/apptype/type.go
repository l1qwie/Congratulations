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

type Notified struct {
	IdSubed       int    `json:"subed_id"`
	IdBirth       int    `json:"birthday_id"`
	NameSubed     string `json:"subed_name"`
	NameBirth     string `json:"birthday_name"`
	NicknameSubed string `json:"subed_nickname"`
	NicknameBirth string `json:"birthday_nickname"`
	Email         string `json:"subed_email"`
	Birthday      string `json:"birthday"`
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
