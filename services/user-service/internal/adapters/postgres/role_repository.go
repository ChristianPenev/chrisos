package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

const rolesNameConstraint = "roles_organization_id_name_key"

type RoleRepository struct {
	db *gorm.DB
}

var _ role.Repository = (*RoleRepository)(nil)

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, ro *role.Role) error {
	rec := roleFromDomain(ro)
	err := r.db.WithContext(ctx).Create(&rec).Error

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation && pgErr.ConstraintName == rolesNameConstraint {
		return role.ErrRoleNameTaken
	}
	return err
}

func (r *RoleRepository) FindByID(ctx context.Context, id role.ID) (*role.Role, error) {
	var rec roleRecord
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", uuid.UUID(id)).Take(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, role.ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return rec.toDomain(), nil
}
