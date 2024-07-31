package tests

import (
	"database/sql"
)

var Con *TestConnection

type TestConnection struct {
	DB *sql.DB
}

func (TC *TestConnection) CreateEmployees() {
	_, err := TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (125, 'Bogdan', 'l1qwie', 'trashemailbyfly@gmail.com', '2003-02-20')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (126, 'Oleja', 'oleja_krut', 'example@gmail.com', '1980-09-20')")
	if err != nil {
		panic(err)
	}
}

func (TC *TestConnection) SubscribeToEmployee() {
	_, err := TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (125, 126, FALSE)")
	if err != nil {
		panic(err)
	}
}

func (TC *TestConnection) DeleteEmployees() {
	_, err := TC.DB.Exec("DELETE FROM Employees")
	if err != nil {
		panic(err)
	}
}

func (TC *TestConnection) DeleteSubscriptions() {
	_, err := TC.DB.Exec("DELETE FROM Subscriptions")
	if err != nil {
		panic(err)
	}
}

func (TC *TestConnection) checkNotifiedInDB(id1, id2 int) bool {
	var count int
	err := TC.DB.QueryRow("SELECT COUNT(*) FROM Subscriptions WHERE subedid = $1 AND subtoid = $2 AND notificated = TRUE", id1, id2).Scan(&count)
	if err != nil {
		panic(err)
	}
	return count > 0
}

func (TC *TestConnection) Create8Employees() {
	_, err := TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (125, 'Bogdan', 'l1qwie', 'trashemailbyfly@gmail.com', '2003-07-30')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (126, 'Oleja', 'oleja_krut', 'example@gmail.com', '1980-08-20')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (127, 'Misha', 'misha228', 'trashfly@ya.ru', '2010-08-10')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (128, 'Katia', 'krasotka_katusha', 'example-example@gmail.com', '1986-12-02')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (129, 'Alesha', 'alexsei', 'me@yahoo.com', '2016-05-28')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (130, 'Nastya', 'princes', 'letmefly@gmail.com', '1970-08-11')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (131, 'Fedor', 'fedor999', 'coolgirl@gmail.com', '1919-08-09')")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (132, 'Matvei', 'cool_juice', 'juice@gmail.com', '1980-08-03')")
	if err != nil {
		panic(err)
	}
}

func (TC *TestConnection) SubEmployeesToEmployees() {
	_, err := TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (127, 131, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (128, 131, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (129, 131, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (127, 126, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (130, 126, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (130, 132, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (127, 132, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (128, 132, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (129, 132, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (130, 125, FALSE)")
	if err != nil {
		panic(err)
	}
	_, err = TC.DB.Exec("INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (128, 125, FALSE)")
	if err != nil {
		panic(err)
	}
}
