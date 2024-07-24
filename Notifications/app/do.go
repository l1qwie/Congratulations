package app

import (
	"Notifications/apptype"
	"database/sql"
	"fmt"
	"net/smtp"
)

var Con *Connection

type Connection struct {
	DB *sql.DB
}

func (c *Connection) findTwoEmployees(id1, id2 int) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Subscriptions WHERE subedid = $1 AND subtoid = $2", id1, id2).Scan(&count)
	if err == nil {
		if count == 0 {
			err = fmt.Errorf("there's nothing about these two employees in the database. The first employee hasn't subscribed to the second")
		}
	}
	return count > 0, err
}

func (c *Connection) updateNotificated() error {
	_, err := c.DB.Exec("UPDATE Subscriptions SET notificated = TRUE WHERE subedid = $1 AND subtoid = $2")
	return err
}

func sendALetter(notified *apptype.Notified) error {
	smtpServer := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", "mya487466@gmail.com", "me_ya_12345566_TestEmail", "smtp.gmail.com")

	// Настройка сообщения
	from := "mya487466@gmail.com"
	to := []string{notified.Email}
	msg := fmt.Sprintf("Hello, %s! I've subdcribed to %s or maybe you could know them us %s for a while and there'll be they Brithday! Their Birthday is going to be on %s This is a message from congratulations service which can notice you to be prepared! See you soon!",
		notified.NameSubed, notified.NameBirth, notified.NicknameBirth, notified.Birthday)

	return smtp.SendMail(smtpServer, auth, from, to, []byte(msg))
}

func Notify(notified *apptype.Notified) (string, error) {
	var answer string
	ok, err := Con.findTwoEmployees(notified.IdSubed, notified.IdBirth)
	if ok && err == nil {
		err = sendALetter(notified)
		if err == nil {
			err = Con.updateNotificated()
		}
	}
	return answer, err
}
