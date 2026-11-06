package app

import "time"

type User struct {
	ID             string
	OrganizationID string
	RoleID         string
	Email          string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
}
