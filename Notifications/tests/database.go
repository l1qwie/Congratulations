package tests

import "database/sql"

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
