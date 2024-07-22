package main

import (
	"time"

	employees "github.com/l1qwie/Congratulations/Employees"
)

func main() {
	time.Sleep(time.Second * 6)
	employees.StartEmployeeTests()
}
