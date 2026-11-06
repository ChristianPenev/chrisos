package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type userRecord struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid"`
	RoleID         uuid.UUID `gorm:"type:uuid"`
	Email          string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
}

func (userRecord) TableName() string { return "users" }

func userFromDomain(u *user.User) userRecord {
	return userRecord{
		ID:             uuid.UUID(u.ID()),
		OrganizationID: uuid.UUID(u.OrganizationID()),
		RoleID:         uuid.UUID(u.RoleID()),
		Email:          u.Email().String(),
		FirstName:      u.FirstName(),
		LastName:       u.LastName(),
		CreatedAt:      u.CreatedAt(),
	}
}

func (r userRecord) toDomain() *user.User {
	return user.Rehydrate(
		user.ID(r.ID),
		organization.ID(r.OrganizationID),
		role.ID(r.RoleID),
		r.Email,
		r.FirstName,
		r.LastName,
		r.CreatedAt,
	)
}
