package app

import (
	"Authorization/apptype"
	"database/sql"
	"fmt"
	"log"
)

var Con *Connection

const (
	nick string = "nickname"
	pass string = "password"
)

type Connection struct {
	DB *sql.DB
}

func (c *Connection) findNickname(a *apptype.Auth) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE nickname = $1 AND password = $2", a.Nickname, a.Password).Scan(&count)
	return count > 0, err
}

func (c *Connection) saveClientIP(a *apptype.Auth, ip string) error {
	_, err := c.DB.Exec("UPDATE Auth SET ip = $1, loggedin = CURRENT_TIMESTAMP WHERE nickname = $2 AND password = $3", ip, a.Nickname, a.Password)
	return err
}

func (c *Connection) checkNickname(nickname string) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE nickname = $1", nickname).Scan(&count)
	return count == 0, err
}

func (c *Connection) savedNewEmployee(auth *apptype.Auth, ip string) (int, error) {
	var (
		id                int
		transactionresult string
	)
	_, err := c.DB.Exec("BEGIN ISOLATION LEVEL REPEATABLE READ")
	if err == nil {
		err = c.DB.QueryRow("SELECT nextval('employeeId')").Scan(&id)
		if err == nil {
			_, err = c.DB.Exec("INSERT INTO Auth (id, nickname, password, ip, loggedin) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)", id, auth.Nickname, auth.Password, ip)
		}
	}
	if err != nil {
		transactionresult = "ROLLBACK"
	} else {
		transactionresult = "COMMIT"
	}
	_, err = c.DB.Exec(transactionresult)
	return id, err
}

func (c *Connection) findEmployeeToDeleteOrChange(id int, nickname, password string) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = $1 AND (nickname = $2 AND password = $3)", id, nickname, password).Scan(&count)
	return count > 0, err
}

func (c *Connection) deleteEmployee(id int) error {
	_, err := c.DB.Exec("DELETE FROM Auth WHERE id = $1", id)
	return err
}

func (c *Connection) checkTimeOut(id int, ip string) (bool, error) {
	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM Auth WHERE id = $1 AND ip = $2 AND (loggedin > (CURRENT_TIMESTAMP - INTERVAL '6 hours'))", id, ip).Scan(&count)
	return count > 0, err
}

func (c *Connection) changeEmployee(id int, column, data string) error {
	query := fmt.Sprintf("UPDATE Auth SET %s = $1 WHERE id = $2", column)
	_, err := c.DB.Exec(query, data, id)
	return err
}

func LogIn(auth *apptype.Auth, clientip string) (string, error) {
	var answer string
	log.Printf("Enter to LogIn() with data:\n	-auth: %v, clientip: %s", auth, clientip)
	ok, err := Con.findNickname(auth)
	if ok {
		log.Print("The nickname and the password was successfuly found")
		log.Print(clientip)
		err = Con.saveClientIP(auth, clientip)
		if err == nil {
			log.Print("Saved client's ip to the database and started timer")
			answer = "You have been authrizated"
		}
	} else {
		log.Print("Didn't find the nickname or the password or neither")
		if err == nil {
			err = fmt.Errorf("incorrect nickname or password")
		}
	}
	return answer, err
}

func SignIn(auth *apptype.Auth, clientip string) (*apptype.SignIn, error) {
	var id int
	log.Printf("Enter to SignIn() with data:\n	-auth: %v, clientip: %s", auth, clientip)
	signin := new(apptype.SignIn)
	ok, err := Con.checkNickname(auth.Nickname)
	if ok {
		id, err = Con.savedNewEmployee(auth, clientip)
		if err == nil {
			signin.Answer = "Welcome! You have successfuly sign in!"
			signin.Id = id
			signin.Nickname = auth.Nickname
			signin.Password = auth.Password
		}
	} else {
		if err == nil {
			err = fmt.Errorf("this nickname is already used. Please come up with diffrent")
		}
	}
	return signin, err
}

func DeleteEmployee(del *apptype.Delete, clientip string) (string, error) {
	var (
		answer string
		timeok bool
	)
	log.Printf("Enter to DeleteEmployee() with data:\n	-del: %v", del)
	ok, err := Con.findEmployeeToDeleteOrChange(del.Id, del.Nickname, del.Passwrod)
	if ok {
		timeok, err = Con.checkTimeOut(del.Id, clientip)
		if timeok {
			err = Con.deleteEmployee(del.Id)
			if err == nil {
				answer = "Employee has been deleted"
			}
		} else {
			if err == nil {
				err = fmt.Errorf("your session has completed. You must authorizate yourself before")
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("couldn't find the employee accroding to the input information. Try to send a diffrent one")
		}
	}
	return answer, err
}

func ChangeEmployee(changes *apptype.Change, clientip string) (string, error) {
	var (
		answer string
		timeok bool
	)
	log.Printf("Enter to ChangeEmployee() with data:\n	-changes: %v, clientip: %s", changes, clientip)
	ok, err := Con.findEmployeeToDeleteOrChange(changes.Id, changes.Nickname, changes.Password)
	if ok {
		if changes.WhichChange == nick || changes.WhichChange == pass {
			timeok, err = Con.checkTimeOut(changes.Id, clientip)
			if timeok {
				err = Con.changeEmployee(changes.Id, changes.WhichChange, changes.ChangeData)
				if err == nil {
					answer = "The employee's data has been just successfuly changed"
				}
			} else {
				if err == nil {
					err = fmt.Errorf("your session has completed. You must authorizate yourself before")
				}
			}
		} else {
			if err == nil {
				err = fmt.Errorf("you must send 'nickname' or 'password' in 'which_change' field")
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("couldn't find the employee accroding to the input information. Try to send a diffrent one")
		}
	}
	return answer, err
}
