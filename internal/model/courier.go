package model

import "github.com/google/uuid"

type Courier struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	WorkTime int       `json:"work_time" db:"work_time"`
	BagSize  int       `json:"bag_size"`
}

type CourierDTO struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Name      string `json:"name"`
	WorkTime  int    `json:"work_time"`
	BagSize   int    `json:"bag_size"`
}
