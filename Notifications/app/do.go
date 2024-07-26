package app

import (
	"Notifications/apptype"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

var Con *Connection

const (
	fromemail string = "cogratulationservice@gmail.com"
	subject   string = "Your Friend's Birthday Is Coming!"
)

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

func (c *Connection) updateNotificated(subedid, subtoid int) error {
	_, err := c.DB.Exec("UPDATE Subscriptions SET notificated = TRUE WHERE subedid = $1 AND subtoid = $2", subedid, subtoid)
	return err
}

func (c *Connection) FindWhoShouldBeNotified() map[*apptype.Employee][]*apptype.Employee {
	var birthdayboys []*apptype.Employee
	members := make(map[*apptype.Employee][]*apptype.Employee)
	rows, err := c.DB.Query(`
		SELECT e.id, e.name, e.nickname, e.birthday 
		FROM Employees e
		JOIN Subscriptions s 
    	ON s.subtoid = e.id 
		WHERE 
    	(
        	(EXTRACT(MONTH FROM CURRENT_DATE) = EXTRACT(MONTH FROM e.birthday) AND EXTRACT(DAY FROM e.birthday) >= EXTRACT(DAY FROM CURRENT_DATE))
        	OR
        	(EXTRACT(MONTH FROM CURRENT_DATE + INTERVAL '30 day') = EXTRACT(MONTH FROM e.birthday) AND EXTRACT(DAY FROM e.birthday) <= EXTRACT(DAY FROM CURRENT_DATE + INTERVAL '30 day'))
        	OR
        	(EXTRACT(MONTH FROM CURRENT_DATE + INTERVAL '30 day') > EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(MONTH FROM e.birthday) BETWEEN EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(MONTH FROM CURRENT_DATE + INTERVAL '30 day'))
    	)
    	AND s.notificated = FALSE
		ORDER BY e.id;`)
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

func sendGomail(to, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromemail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromemail, "ycuw acml gnor qcir")
	err := d.DialAndSend(m)
	if err != nil {
		log.Print(err)
	}
	return err
}

func sendALetter(notified *apptype.Notified) error {
	var err error
	birthday, err := time.Parse("2006-01-02", notified.Birthday)
	if err == nil {
		today := time.Now()
		age := today.Year() - birthday.Year()
		nextBirthday := time.Date(today.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)
		if today.After(nextBirthday) {
			nextBirthday = time.Date(today.Year()+1, birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)
		}
		if today.YearDay() < birthday.YearDay() {
			age--
		}
		duration := nextBirthday.Sub(today)
		days := int(duration.Hours() / 24)
		message := fmt.Sprintf(`<b>Hello, %s "%s"!</b>`, notified.NameSubed, notified.NicknameSubed) + "\n" +
			fmt.Sprintf(`You've subscribed to <b>%s</b> for a while, and it will be their <b>Birthday</b> soon! His/Her age turns <b>%d</b>! Their were born in <b>%s</b> and you have only <b>%d</b> days before it'll come!`,
				notified.NameBirth, age, notified.Birthday, days) + "\n" + "This is a message from the congratulations service, which can notify you to be prepared! Good luck!"
		err = sendGomail(notified.Email, message)
	}

	return err
}

func prepareLetters(members map[*apptype.Employee][]*apptype.Employee) {
	for key, value := range members {
		for _, empl := range value {
			birthday, err := time.Parse("2006-01-02", key.Birthday)
			if err == nil {
				today := time.Now()
				age := today.Year() - birthday.Year()
				nextBirthday := time.Date(today.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)
				if today.After(nextBirthday) {
					nextBirthday = time.Date(today.Year()+1, birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)
				}
				if today.YearDay() < birthday.YearDay() {
					age--
				}
				duration := nextBirthday.Sub(today)
				days := int(duration.Hours() / 24)
				message := fmt.Sprintf(`<b>Hello, %s "%s"!</b>`, empl.Name, empl.Nickname) + "\n" +
					fmt.Sprintf(`You've subscribed to <b>%s</b> for a while, and it will be their <b>Birthday</b> soon! His/Her age turns <b>%d</b>! Their were born in <b>%s</b> and you have only <b>%d</b> days before it'll come!`,
						key.Name, age, key.Birthday, days) + "\n" + "This is a message from the congratulations service, which can notify you to be prepared! Good luck!"
				sendGomail(empl.Email, message)
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
			err = Con.updateNotificated(notified.IdSubed, notified.IdBirth)
			if err == nil {
				answer = "The employee was notified"
			}
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
