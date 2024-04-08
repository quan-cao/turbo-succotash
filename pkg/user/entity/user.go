package entity

import "time"

type User struct {
	Id        int
	Isid      string
	Role      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
