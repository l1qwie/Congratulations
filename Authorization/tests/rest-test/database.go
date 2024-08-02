package resttest

import (
	"Authorization/apptype"
	"database/sql"
	"log"
)

var Con *TestConnection

type TestConnection struct {
	DB *sql.DB
}

func (c *TestConnection) checkLoggedIn(a *apptype.Auth) bool {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE nickname = $1 AND password = $2 AND loggedin = CURRENT_TIMESTAMP",
		a.Nickname, a.Password).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (c *TestConnection) CreateEmployee() {
	_, err := c.DB.Exec("INSERT INTO Auth (id, nickname, password, ip, loggedin) VALUES (nextval('employeeId'), 'mama-miya', '12345678', '::1', CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	log.Print("Create an employee")
}

func (c *TestConnection) CreateEmployeeKafka() {
	_, err := c.DB.Exec("INSERT INTO Employees (id, nickname) VALUES (1, 'mama-miya')")
	if err != nil {
		panic(err)
	}
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

func (c *TestConnection) ResetSequence() {
	_, err := c.DB.Exec("SELECT setval('employeeId', 1, false)")
	if err != nil {
		panic(err)
	}
	log.Print("Sequence 'employeeId' has been reseted")
}

func (c *TestConnection) ckeckNewEmployee(signin *apptype.SignIn) bool {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = $1 AND nickname = $2 AND password = $3 AND ip = '::1'", signin.Id, signin.Nickname, signin.Password).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (c *TestConnection) ckeckDeleteEmployee(id int) bool {
	var count int
	err := c.DB.QueryRow("SELECT COunT(*) FROM Auth WHERE id = $1", id).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count == 0
}

func (c *TestConnection) checkChangedNicknameEmployee(nickname string) bool {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = 1 AND (nickname = $1 AND password = '12345678')", nickname).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (c *TestConnection) checkChangedPasswordEmployee(password string) bool {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = 1 AND (nickname = 'mama-miya' AND password = $1)", password).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (c *TestConnection) CkeckAddedOrUpdatedEmployeeKafka(empl *apptype.KafkaEmployee) (bool, bool) {
	var count, count1 int

	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = $1 AND nickname = $2", empl.Id, empl.Nickname).Scan(&count)
	if err != nil {
		panic(err)
	}
	err = c.DB.QueryRow("SELECT COUNT(*) FROM Employees WHERE id = $1 AND name = $2 AND nickname = $3 AND email = $4 AND birthday = $5",
		empl.Id, empl.Name, empl.Nickname, empl.Email, empl.Birthday).Scan(&count1)
	if err != nil {
		panic(err)
	}
	return count == 1, count1 == 1
}

func (c *TestConnection) CheckDeletedEmployeeKafka(empl *apptype.KafkaEmployee) (bool, bool, bool) {
	var count int

	auth, employess := c.CkeckAddedOrUpdatedEmployeeKafka(empl)
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Subscriptions WHERE (subtoid = $1) OR (subedid = $1)", empl.Id).Scan(&count)
	if err != nil {
		panic(err)
	}
	return auth, employess, count == 0
}

func (c *TestConnection) DeleteDONOTUSE() {
	_, err := c.DB.Exec("DELETE FROM DONOTUSE")
	if err != nil {
		panic(err)
	}
}
