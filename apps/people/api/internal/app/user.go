package app

import "time"

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
}
