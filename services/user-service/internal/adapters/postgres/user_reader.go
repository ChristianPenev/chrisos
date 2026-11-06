package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
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

	return query.UserView{
		ID:             rec.ID.String(),
		OrganizationID: rec.OrganizationID.String(),
		RoleID:         rec.RoleID.String(),
		Email:          rec.Email,
		FirstName:      rec.FirstName,
		LastName:       rec.LastName,
		CreatedAt:      rec.CreatedAt.UTC(),
	}, nil
}
