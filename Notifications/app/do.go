package app

import (
	"Notifications/apptype"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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

func (c *Connection) FindWhoShouldBeNotified() map[*apptype.Employee][]*apptype.Employee {
	var birthdayboys []*apptype.Employee
	members := make(map[*apptype.Employee][]*apptype.Employee)
	rows, err := c.DB.Query(`
		SELECT e.id, e.name, e.nickname, e.birthday FROM Employees e
		JOIN Subscriptions s ON s.subtoid = e.id
		WHERE s.notificated = FALSE
		ORDER BY e.id`)
	if err == nil {
		defer rows.Close()
		for rows.Next() && err == nil {
			employee := new(apptype.Employee)
			err = rows.Scan(&employee.Id, &employee.Name, &employee.Nickname, &employee.Birthday)
			if err == nil {
				birthdayboys = append(birthdayboys, employee)
			}
		}
	}
	if err == nil {
		for _, birthdayboy := range birthdayboys {
			rows, err = c.DB.Query(`
				SELECT e.id, e.name, e.nickname, e.email FROM Employees e
				JOIN Subscriptions s ON s.subedid = e.id
				WHERE s.subtoid = $1 AND s.notificated = FALSE`, birthdayboy.Id)
			if err == nil {
				var congratulators []*apptype.Employee
				defer rows.Close()
				for rows.Next() && err == nil {
					employee := new(apptype.Employee)
					err = rows.Scan(&employee.Id, &employee.Name, &employee.Nickname, &employee.Email)
					if err == nil {
						congratulators = append(congratulators, employee)
					}
				}
				members[birthdayboy] = congratulators
			}
		}
	}
	return members
}

func sendALetter(notified *apptype.Notified) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err == nil {
		svc := ses.New(sess)
		message := fmt.Sprintf("Hello, %s %s! You've subscribed to %s or maybe you know them as %s for a while, and it will be their Birthday soon! Their Birthday is going to be on %s. This is a message from the congratulations service, which can notify you to be prepared! Good luck!",
			notified.NameSubed, notified.NicknameSubed, notified.NameBirth, notified.NicknameBirth, notified.Birthday)
		err = sendEmail("mya487466@gmail.com", message, "Birthday Reminder!", []string{notified.Email}, svc)
	}
	return err
}

func sendEmail(from, body, subject string, to []string, svc *ses.SES) error {
	msg := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(to),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	_, err := svc.SendEmail(msg)
	if err != nil {
		log.Print(err)
	}
	return err
}

func prepareLetters(members map[*apptype.Employee][]*apptype.Employee) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err == nil {
		from := "mya487466@gmail.com"
		subject := "Birthday Reminder!"
		svc := ses.New(sess)
		for key, value := range members {
			for _, empl := range value {
				message := fmt.Sprintf("Hello, %s %s! You've subscribed to %s or maybe you know them as %s for a while, and it will be their Birthday soon! Their Birthday is going to be on %s. This is a message from the congratulations service, which can notify you to be prepared! Good luck!",
					empl.Name, empl.Nickname, key.Name, key.Nickname, key.Birthday)
				sendEmail(from, message, subject, []string{empl.Email}, svc)
			}

		}
	}
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

func Guardian() {
	for {
		members := Con.FindWhoShouldBeNotified()
		if len(members) > 0 {
			prepareLetters(members)
		}
		time.Sleep(time.Second)
	}
}
