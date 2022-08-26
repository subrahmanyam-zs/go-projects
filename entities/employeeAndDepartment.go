package entities

import "github.com/google/uuid"

type EmployeeAndDepartment struct {
	ID     uuid.UUID
	Name   string
	Dob    string
	City   string
	Majors string
	Dept   Department
}
