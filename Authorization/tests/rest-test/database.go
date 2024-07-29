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
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE nickname = $1 AND password = $2 AND ip = '::1'",
		a.Nickname, a.Password).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (c *TestConnection) createEmployee() {
	_, err := c.DB.Exec("INSERT INTO Auth (id, nickname, password, ip, loggedin) VALUES (nextval('employeeId'), 'mama-miya', '12345678', '::1', CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	log.Print("Create an employee")
}

func (c *TestConnection) deleteEmployees() {
	_, err := c.DB.Exec("DELETE FROM Auth")
	if err != nil {
		panic(err)
	}
	log.Print("All data has been deleted from table Auth")
}

func (c *TestConnection) resetSequence() {
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
