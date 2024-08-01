package resttest

import (
	"database/sql"
	"log"
)

var Con *TestConnection

type TestConnection struct {
	DB *sql.DB
}

func (c *TestConnection) ResetSequence() {
	_, err := c.DB.Exec("SELECT setval('employeeId', 1, false)")
	if err != nil {
		panic(err)
	}
	log.Print("Sequence 'employeeId' has been reseted")
}

func (c *TestConnection) DeleteEmployees() {
	_, err := c.DB.Exec("DELETE FROM Employees")
	if err != nil {
		panic(err)
	}
	_, err = c.DB.Exec("DELETE FROM Auth")
	if err != nil {
		panic(err)
	}
	log.Print("All data has been deleted from table Auth")
}

func (c *TestConnection) CreateEmployee() {
	_, err := c.DB.Exec("INSERT INTO Auth (id, nickname, password, ip, loggedin) VALUES (nextval('employeeId'), 'mama-miya', '12345678', '::1', CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	_, err = c.DB.Exec("INSERT INTO Auth (id, nickname, password, ip, loggedin) VALUES (nextval('employeeId'), 'shuman', '12358', '::1', CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	_, err = c.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (1, '', 'mama-miya', 'ekkk@gmail.com', '2009-02-05')")
	if err != nil {
		panic(err)
	}
	_, err = c.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (2, '', 'shuman', 'example@yahoo.com', '1977-09-09')")
	if err != nil {
		panic(err)
	}
	log.Print("Created 4 employees")
}

func (c *TestConnection) DeleteSubscribers() {
	_, err := c.DB.Exec("DELETE FROM Subscriptions")
	if err != nil {
		panic(err)
	}
}

func (c *TestConnection) CheckSubedEmployees(id, id1 int) bool {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Subscriptions WHERE subedid = $1 AND subtoid = $2", id, id1).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count == 1
}

func (c *TestConnection) SubscribeEmployeeToEmployee() {
	_, err := c.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid) VALUES (1, 2)")
	if err != nil {
		panic(err)
	}
}
