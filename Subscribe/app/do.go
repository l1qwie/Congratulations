package app

import (
	producer "Subscribe/api/kafka/poducer"
	"database/sql"
	"fmt"
)

var Con *Connection

type Connection struct {
	DB *sql.DB
}

const (
	Sub   string = "sub"
	Unsub string = "unsub"
)

func (c *Connection) FindEmploee(id int) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Employees WHERE id = $1", id).Scan(&count)
	return count == 1, err
}

func (c *Connection) CheckSubStatus(subscriber, subto int) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Subscriptions WHERE subedid = $1 AND subtoid = $2", subscriber, subto).Scan(&count)
	return count == 0, err
}

func (c *Connection) SubEmployeeToEmployee(subscriber, subto int) error {
	_, err := c.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid) VALUES ($1, $2)", subscriber, subto)
	return err
}

func (c *Connection) UnsubEmployeeFromEmployee(subscriber, subto int) error {
	_, err := c.DB.Exec("DELETE FROM Subscriptions WHERE subedid = $1 AND subtoid = $2", subscriber, subto)
	return err
}

func SubscribeAction(subscriber, subtoid int, whatdo string) (string, error) {
	var answer string
	ok, err := Con.FindEmploee(subscriber)
	if ok {
		ok, err = Con.FindEmploee(subtoid)
		if ok {
			ok, err = Con.CheckSubStatus(subscriber, subtoid)

			if ok && whatdo == Sub {
				err = Con.SubEmployeeToEmployee(subscriber, subtoid)
				if err == nil {
					answer = fmt.Sprintf("You've just subscribed to %d", subtoid)
					producer.TellChanges(subscriber, subtoid, Sub)
				}

			} else if !ok && whatdo == Unsub {
				err = Con.UnsubEmployeeFromEmployee(subscriber, subtoid)
				if err == nil {
					answer = fmt.Sprintf("You've just unsubscribed from %d", subtoid)
					producer.TellChanges(subscriber, subtoid, Unsub)
				}

			} else if !ok && whatdo == Sub {
				if err == nil {
					err = fmt.Errorf("this employee: %d has already subscribed to this: %d", subscriber, subtoid)
				}

			} else if ok && whatdo == Unsub {
				if err == nil {
					err = fmt.Errorf("this employee: %d doesn't subscribe to this: %d", subscriber, subtoid)
				}
			}
		} else {
			if err == nil {
				err = fmt.Errorf("(un)subscribe_to's id isn't found")
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("subscriber's id isn't found")
		}
	}
	return answer, err
}
