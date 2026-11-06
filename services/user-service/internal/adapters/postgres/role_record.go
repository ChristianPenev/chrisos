package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

type roleRecord struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid"`
	Name           string
	Owner          bool
	CreatedAt      time.Time
	Permissions    []rolePermissionRecord `gorm:"foreignKey:RoleID"`
}

func (roleRecord) TableName() string { return "roles" }

type rolePermissionRecord struct {
	RoleID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Permission string    `gorm:"primaryKey"`
}

func (rolePermissionRecord) TableName() string { return "role_permissions" }

func roleFromDomain(r *role.Role) roleRecord {
	rec := roleRecord{
		ID:             uuid.UUID(r.ID()),
		OrganizationID: uuid.UUID(r.OrganizationID()),
		Name:           r.Name(),
		Owner:          r.IsOwner(),
		CreatedAt:      r.CreatedAt(),
	}
	if r.IsOwner() {
		return rec
	}
	for _, p := range r.Permissions() {
		rec.Permissions = append(rec.Permissions, rolePermissionRecord{RoleID: rec.ID, Permission: p.String()})
	}
	return rec
}

func (r roleRecord) toDomain() *role.Role {
	perms := make([]string, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = p.Permission
	}
	return role.Rehydrate(role.ID(r.ID), organization.ID(r.OrganizationID), r.Name, perms, r.Owner, r.CreatedAt)
}
