package app

import "time"

type Role struct {
	ID             string
	OrganizationID string
	Name           string
	IsOwner        bool
	Permissions    []string
	CreatedAt      time.Time
}
