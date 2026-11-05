package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type userRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
}

func (userRecord) TableName() string { return "users" }

func fromDomain(u *user.User) userRecord {
	return userRecord{
		ID:        uuid.UUID(u.ID()),
		Email:     u.Email().String(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
		CreatedAt: u.CreatedAt(),
	}
}

func (r userRecord) toDomain() *user.User {
	return user.Rehydrate(user.ID(r.ID), r.Email, r.FirstName, r.LastName, r.CreatedAt)
}
