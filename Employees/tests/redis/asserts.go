package redis

import (
	"context"
	"fmt"
	"reflect"

	"github.com/l1qwie/Congratulations/Employees/apptype"
)

// Проверяет удаленного работника из хэш-таблицы
func (TRCL *TestRedClient) CheckDeletedEmployee(employeeid string) bool {
	exists, err := TRCL.Cl.HExists(context.Background(), "employeeid", employeeid).Result()
	if err != nil {
		panic(err)
	}
	return exists
}

// Проверяет на нового работника или на обновленного (по сути тоже нового)
func (TRCL *TestRedClient) CheckUpdatedOrNewEmployee(req *apptype.Employee) bool {
	key := fmt.Sprintf("employeeid:%d", req.Id)
	empl := new(apptype.Employee)
	fields, err := TRCL.Cl.HGetAll(context.Background(), key).Result()
	if err == nil {
		empl = &apptype.Employee{
			Id:       req.Id,
			Name:     fields["name"],
			Nickname: fields["nickname"],
			Email:    fields["email"],
			Birthday: fields["birthday"],
		}
	}
	return reflect.DeepEqual(req, empl)
}

// Проверяет на сабнутого работника или ансабнутого
func (TRCL *TestRedClient) CheckSubToEmployee(subedid, subtoid int) bool {
	key := fmt.Sprintf("subscriptions:%d", subedid)
	exists, err := TRCL.Cl.SIsMember(context.Background(), key, subtoid).Result()
	if err != nil {
		panic(err)
	}
	return exists
}
