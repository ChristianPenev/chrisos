package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	pgshared "github.com/ChristianPenev/chrisos/pkg/postgres"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type UserReader struct {
	db *gorm.DB
}

var _ query.UserReader = (*UserReader)(nil)

func NewUserReader(db *gorm.DB) *UserReader {
	return &UserReader{db: db}
}

func (r *UserReader) UserByID(ctx context.Context, id user.ID) (query.UserView, error) {
	var rec userRecord
	err := r.db.WithContext(ctx).Where("id = ?", uuid.UUID(id)).Take(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.UserView{}, user.ErrUserNotFound
	}
	if err != nil {
		return query.UserView{}, err
	}

	return rec.toView(), nil
}

func (r *UserReader) UsersByOrganization(ctx context.Context, organizationID organization.ID, req pagination.Request) (pagination.Page[query.UserView], error) {
	recs, err := pgshared.Keyset(
		r.db.WithContext(ctx).Where("organization_id = ?", uuid.UUID(organizationID)),
		req,
		"id",
		func(rec userRecord) uuid.UUID { return rec.ID },
	)
	if err != nil {
		return pagination.Page[query.UserView]{}, err
	}

	views := make([]query.UserView, len(recs.Items))
	for i, rec := range recs.Items {
		views[i] = rec.toView()
	}
	return pagination.Page[query.UserView]{Items: views, NextCursor: recs.NextCursor}, nil
}

func (r userRecord) toView() query.UserView {
	return query.UserView{
		ID:             r.ID.String(),
		OrganizationID: r.OrganizationID.String(),
		RoleID:         r.RoleID.String(),
		Email:          r.Email,
		FirstName:      r.FirstName,
		LastName:       r.LastName,
		CreatedAt:      r.CreatedAt.UTC(),
	}
}
