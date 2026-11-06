package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

type organizationRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	CreatedAt time.Time
}

func (organizationRecord) TableName() string { return "organizations" }

func organizationFromDomain(o *organization.Organization) organizationRecord {
	return organizationRecord{
		ID:        uuid.UUID(o.ID()),
		Name:      o.Name(),
		CreatedAt: o.CreatedAt(),
	}
}

func (r organizationRecord) toDomain() *organization.Organization {
	return organization.Rehydrate(organization.ID(r.ID), r.Name, r.CreatedAt)
}
