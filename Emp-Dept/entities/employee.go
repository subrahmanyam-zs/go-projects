package entities

import "github.com/google/uuid"

type Employee struct {
	ID     uuid.UUID `json:"ID,omitempty"`
	Name   string    `json:"Name,omitempty"`
	Dob    string    `json:"Dob,omitempty"`
	City   string    `json:"City,omitempty"`
	Majors string    `json:"majors,omitempty"`
	DeptID int       `json:"DeptID,omitempty"`
}
