package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/l1qwie/Congratulations/Employees/apptype"
)

type TestRedClient struct {
	Cl *redis.Client
}

func AddClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	return client, err
}

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
	}
}

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
}

func (TRCL *TestRedClient) DeleteEmployees() {
	err := TRCL.Cl.FlushAll(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

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
