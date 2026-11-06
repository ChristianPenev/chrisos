package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

type OrganizationRepository struct {
	db *gorm.DB
}

var _ organization.Repository = (*OrganizationRepository)(nil)

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, o *organization.Organization) error {
	rec := organizationFromDomain(o)
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id organization.ID) (*organization.Organization, error) {
	var rec organizationRecord
	err := r.db.WithContext(ctx).Where("id = ?", uuid.UUID(id)).Take(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, organization.ErrOrganizationNotFound
	}
	if err != nil {
		return nil, err
	}
	return rec.toDomain(), nil
}
