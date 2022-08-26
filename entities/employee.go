package entities

import "github.com/google/uuid"

type Employee struct {
	ID     uuid.UUID
	Name   string
	Dob    string
	City   string
	Majors string
	DId    int
}
