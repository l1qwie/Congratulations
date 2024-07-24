package app

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"Employees/api/kafka/producer"
	"Employees/apptype"

	"github.com/go-redis/redis/v8"
)

const deflimit int = 15

type RedClient struct {
	Cl *redis.Client
}

var Client *RedClient

// Получает данные из хэш-таблицы по определенному employeeid
func (RCL *RedClient) getEmployeeByID(subid int) (*apptype.Employee, error) {
	key := fmt.Sprintf("employeeid:%d", subid)
	empl := new(apptype.Employee)
	fields, err := RCL.Cl.HGetAll(context.Background(), key).Result()
	if err == nil {
		empl = &apptype.Employee{
			Id:       subid,
			Name:     fields["name"],
			Nickname: fields["nickname"],
			Email:    fields["email"],
			Birthday: fields["birthday"],
		}
	}
	return empl, err
}

// Вызывает функции, которые вынимает нужные данные из множества и хэш таблицы, а так же формирует слайс.
// Формирует слайс из данных подписсчиков одного работника
func (RCL *RedClient) selectSubEmployees(employeeid, limit int) ([]*apptype.Employee, error) {
	var employees []*apptype.Employee
	log.Printf("Got into selectSubEmployees with params: employeeid: %d, limit: %d", employeeid, limit)
	key := fmt.Sprintf("subscriptions:%d", employeeid)
	subs, err := RCL.Cl.SMembers(context.Background(), key).Result()
	if err == nil {
		log.Print("It has succsessfuly just pulled all subids from redis")
		subid := 0
		log.Print(len(subs))
		for i := 0; i < len(subs) && err == nil; i++ {
			subid, err = strconv.Atoi(subs[i])
			if err == nil {
				log.Print("It has successfuly converted string to int")
				empl, err := RCL.getEmployeeByID(subid)
				if err == nil {
					log.Print("Successfuly got an employee by its ID")
					employees = append(employees, empl)
				}
			}
		}
	}
	log.Print("Get out of selectSubEmployees")
	return employees, err
}

// Получает все id работников из множества
func (RCL *RedClient) getAllEmployeeIDs() ([]string, error) {
	key := "employee_ids"
	ids, err := RCL.Cl.SMembers(context.Background(), key).Result()
	return ids, err
}

// Вызывает функции, которые вынимает нужные данные из множества и хэш таблицы, а так же формирует слайс
// Формирует слайс из всевозможных данных о всех работниках
func (RCL *RedClient) selectEmployees(limit int) ([]*apptype.Employee, error) {
	var employees []*apptype.Employee
	log.Printf("Got into selectEmployees with a param: limit: %d", limit)
	ids, err := RCL.getAllEmployeeIDs()
	if err == nil {
		log.Print("Got all employees IDs")
		if len(ids) > limit {
			log.Print("The length of ids is bigger than limit")
			ids = ids[:limit]
		}
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				log.Print("Successfuly converted string to int")
				empl, err := RCL.getEmployeeByID(id)
				if err == nil {
					log.Print("Successfuly got an employee by its ID")
					employees = append(employees, empl)
				}
			}
		}
	}
	log.Print("Got out of selectEmployees")
	return employees, err
}

// Добавляет работника в хэш-таблицу и также добавляет id работника в множество
func (RCL *RedClient) NewEmpl(empl *apptype.Employee) error {
	key := fmt.Sprintf("employeeid:%d", empl.Id)
	_, err := RCL.Cl.HSet(context.Background(), key, map[string]interface{}{
		"name":     empl.Name,
		"nickname": empl.Nickname,
		"email":    empl.Email,
		"birthday": empl.Birthday,
	}).Result()
	if err == nil {
		idKey := "employee_ids"
		_, err = RCL.Cl.SAdd(context.Background(), idKey, empl.Id).Result()
	}
	return err
}

// Удлаяет работника из хэш-таблицы
func (RCL *RedClient) DeleleEmpl(id int) error {
	key := fmt.Sprintf("employeeid:%d", id)
	err := RCL.Cl.Del(context.Background(), key).Err()
	if err != nil {
		log.Print(err)
	}
	return err
}

// Обновляет данные работника. По сути - вызывает функцию удаления и функцию добавления нового работника
func (RCL *RedClient) UpdEmpl(empl *apptype.Employee, id int) error {
	err := RCL.DeleleEmpl(id)
	if err == nil {
		err = RCL.NewEmpl(empl)
	}
	return err
}

// Добавляет значение в множество (подписывает)
func (RCL *RedClient) AddSub(employeeid, subtoId int) error {
	key := fmt.Sprintf("subscriptions:%d", employeeid)
	_, err := RCL.Cl.SAdd(context.Background(), key, subtoId).Result()
	return err
}

// Убирает значение из множества (отписывает)
func (RCL *RedClient) UnSub(employeeid, unsubtoId int) error {
	key := fmt.Sprintf("subscriptions:%d", employeeid)
	_, err := RCL.Cl.SRem(context.Background(), key, unsubtoId).Result()
	return err
}

// Делает поиск нужного employeeID в множестве
func (RCL *RedClient) findEmployee(id int) (bool, error) {
	found, err := RCL.Cl.SIsMember(context.Background(), "employee_ids", id).Result()
	return found, err
}

// Делает предварительные проверки, а конкретнее: на существование employeeID в бд и значение переменной whatdo и перенавправляет на нужную функцию
func (RCL *RedClient) updateEmployee(empl *apptype.Employee, whatdo, diffrentemplid string) error {
	var (
		err     error
		id      int
		ok, ok2 bool
	)
	ok = true
	ok2 = true
	log.Printf("Got in updateEmployee with params: empl: %v, whatdo: %s, diffrentemplid: %s", *empl, whatdo, diffrentemplid)
	ok, err = RCL.findEmployee(empl.Id)
	if whatdo == "update" || whatdo == "sub" || whatdo == "unsub" {
		id, err = strconv.Atoi(diffrentemplid)
		if err == nil {
			log.Print("Has successfuly converted string to int")
			ok2, err = RCL.findEmployee(id)
			log.Print(ok2, err)
		}
	}
	if whatdo == "update" && !ok || ok && ok2 {
		if whatdo == "sub" {
			log.Print(`Whatdo is "delete"`)
			err = RCL.AddSub(empl.Id, id)
		} else if whatdo == "unsub" {
			log.Print(`Whatdo is "delete"`)
			err = RCL.UnSub(empl.Id, id)
		} else if whatdo == "delete" {
			log.Print(`Whatdo is "delete"`)
			err = RCL.DeleleEmpl(empl.Id)
		} else if whatdo == "update" {
			log.Print(`Whatdo is "update"`)
			err = RCL.UpdEmpl(empl, id)
		}
	} else if whatdo == "new" && !ok {
		log.Print(`Whatdo is "new"`)
		err = RCL.NewEmpl(empl)
	} else {
		if err == nil {
			if whatdo == "update" {
				err = fmt.Errorf("you must send an uncreated new employee.Id")
			} else {
				err = fmt.Errorf("wasn't able to find recieved employee ID. You might try to send a diffrent one")
			}
		}
	}

	log.Print("Got out of updateEmployee")
	return err
}

// Мейн функция для get endPoint'a
func GetEmployees(id, limit int) ([]*apptype.Employee, error) {
	var (
		employees []*apptype.Employee
		err       error
	)
	log.Print("Got into GetEmployees func in the bussines logic")
	if id != 0 {
		log.Print("The recieved id isn't 0")
		if limit != 0 {
			log.Print("The recieved limit isn't 0")
			employees, err = Client.selectSubEmployees(id, limit)
		} else {
			log.Print("The recieved limit is 0")
			employees, err = Client.selectSubEmployees(id, deflimit)
		}
	} else {
		log.Print("The recieved id is 0")
		if limit != 0 {
			log.Print("The recieved limit isn't 0")
			employees, err = Client.selectEmployees(limit)
		} else {
			log.Print("The recieved limit is 0")
			employees, err = Client.selectEmployees(deflimit)
		}
	}
	log.Print("Got out of GetEmployees func in the bussines logic")
	return employees, err
}

// Мейн функция для post endPoint'a
func UpdateEmployees(empl *apptype.Employee, whatdo, diffrentemplid string) (string, error) {
	var (
		answer string
		err    error
	)
	log.Print("Got into UpdateEmployees func in the bussines logic")
	err = Client.updateEmployee(empl, whatdo, diffrentemplid)
	if err == nil {
		log.Print("Has successfuly updated redis")
		answer = "The employee has been updated"
		// Не важно, были ли данные в diffrentemplid так как в противном случае - просто будет передан дефолт значение int
		id, _ := strconv.Atoi(diffrentemplid)
		producer.TellChanges(empl, whatdo, id)
	}
	log.Print("Got out of UpdateEmployees func in the bussines logic")
	return answer, err
}
