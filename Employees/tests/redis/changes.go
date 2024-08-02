package redis

import (
	"context"
	"fmt"

	"Employees/apptype"

	"github.com/go-redis/redis/v8"
)

type TestRedClient struct {
	Cl *redis.Client
}

// Добавление редис-клиента test-part
func AddClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "congratulations-redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	return client, err
}

// Создает 1 работника в редис
func (TRCL *TestRedClient) createEmpl(empl *apptype.Employee) {
	key := fmt.Sprintf("employeeid:%d", empl.Id)
	_, err := TRCL.Cl.HSet(context.Background(), key, map[string]interface{}{
		"name":     empl.Name,
		"nickname": empl.Nickname,
		"email":    empl.Email,
		"birthday": empl.Birthday,
	}).Result()
	if err == nil {
		idKey := "employee_ids"
		_, err = TRCL.Cl.SAdd(context.Background(), idKey, empl.Id).Result()
		if err != nil {
			panic(err)
		}
	}
}

// Создает данные работников и отправляет их в функцию, которая добавляет их в бд
func (TRCL *TestRedClient) InitEmployees() {
	// The first employee
	empl := &apptype.Employee{
		Id:       111,
		Name:     "John",
		Nickname: "John1234",
		Email:    "example@ya.ru",
		Birthday: "1990-03-15",
	}
	TRCL.createEmpl(empl)
	// The second employee
	empl = &apptype.Employee{
		Id:       122,
		Name:     "Moisha",
		Nickname: "l1qwie",
		Email:    "example@gmail.com",
		Birthday: "2003-02-20",
	}
	TRCL.createEmpl(empl)
	// The third employee
	empl = &apptype.Employee{
		Id:       172,
		Name:     "Arkadii",
		Nickname: "nagibator_2000",
		Email:    "example@yahoo.com",
		Birthday: "2016-09-15",
	}
	TRCL.createEmpl(empl)
	// The fourth employee
	empl = &apptype.Employee{
		Id:       199,
		Name:     "Masha",
		Nickname: "krasotka-mashka",
		Email:    "example@.dark.net",
		Birthday: "1960-01-18",
	}
	TRCL.createEmpl(empl)
}

// Удаляет вообще все из бд (редис)
func (TRCL *TestRedClient) DeleteEmployees() {
	err := TRCL.Cl.FlushAll(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

// Подписывает отдного работника на трех других
func (TRCL *TestRedClient) SubscribeEmployees() {
	key := fmt.Sprintf("subscriptions:%d", 111)
	_, err := TRCL.Cl.SAdd(context.Background(), key, "111").Result()
	if err != nil {
		panic(err)
	}
	_, err = TRCL.Cl.SAdd(context.Background(), key, "122").Result()
	if err != nil {
		panic(err)
	}
	_, err = TRCL.Cl.SAdd(context.Background(), key, "172").Result()
	if err != nil {
		panic(err)
	}
}
