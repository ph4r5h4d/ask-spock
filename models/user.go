package models

import "time"

type User struct {
	Id        int64
	Username  string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
