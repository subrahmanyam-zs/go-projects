package entities

import "github.com/google/uuid"

type Employee struct {
	Id     uuid.UUID
	Name   string
	Dob    string
	City   string
	Majors string
	DId    int
}
