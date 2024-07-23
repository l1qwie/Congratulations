package kafka

import (
	"Employees/api/kafka/producer"
	"Employees/apptype"
	"Employees/tests/redis"
	"log"
	"time"
)

func testNew(TRCL *redis.TestRedClient) {
	log.Print("testNew has been started")
	employee := &apptype.Employee{
		Id:       168,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer.TellChanges(employee, "new", 0)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckUpdatedOrNewEmployee(employee, false, "") {
		panic("The new employee wasn't added")
	}
	log.Print("testNew has been finished")
}

func testUpdate(TRCL *redis.TestRedClient) {
	log.Print("testUpdate has been started")
	employee := &apptype.Employee{
		Id:       168,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer.TellChanges(employee, "update", 111)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckUpdatedOrNewEmployee(employee, true, "111") {
		panic("The employee wasn't updated")
	}
	log.Print("testUpdate has been finished")
}

func testDelete(TRCL *redis.TestRedClient) {
	log.Print("testDelete has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer.TellChanges(employee, "delete", 0)
	time.Sleep(time.Second / 10)
	if TRCL.CheckDeletedEmployee("111") {
		panic("The employee wasn't deleted")
	}
	log.Print("testDelete has been finished")
}

func testSub(TRCL *redis.TestRedClient) {
	log.Print("testSub has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer.TellChanges(employee, "sub", 199)
	time.Sleep(time.Second / 10)
	if !TRCL.CheckSubToEmployee(employee.Id, 199) {
		panic("The employee wasn't subed to another employee")
	}
	log.Print("testSub has been finished")
}

func testUnsub(TRCL *redis.TestRedClient) {
	log.Print("testUnsub has been started")
	employee := &apptype.Employee{
		Id:       111,
		Name:     "Zick",
		Nickname: "zick_grut",
		Email:    "zick@gmail.com",
		Birthday: "2000-02-20",
	}
	producer.TellChanges(employee, "unsub", 111)
	time.Sleep(time.Second / 10)
	if TRCL.CheckSubToEmployee(employee.Id, 111) {
		panic("The employee wasn't unsubed to another employee")
	}
	log.Print("testUnsub has been finished")
}

func TestConsumer() {
	log.Print("TestConsumer has been started")
	var err error
	TRCL := new(redis.TestRedClient)
	TRCL.Cl, err = redis.AddClient()
	if err != nil {
		panic(err)
	}

	TRCL.InitEmployees()
	TRCL.SubscribeEmployees()
	defer TRCL.DeleteEmployees()

	time.Sleep(time.Second / 10)

	testNew(TRCL)
	testUpdate(TRCL)
	testDelete(TRCL)
	testSub(TRCL)
	testUnsub(TRCL)
	log.Print("TestConsumer has been finished")
}
