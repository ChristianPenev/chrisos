package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
)

type RoleReader struct {
	db *gorm.DB
}

var _ query.RoleReader = (*RoleReader)(nil)

func NewRoleReader(db *gorm.DB) *RoleReader {
	return &RoleReader{db: db}
}

func (r *RoleReader) RolesByIDs(ctx context.Context, ids []role.ID) ([]query.RoleView, error) {
	keys := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		keys[i] = uuid.UUID(id)
	}

	var recs []roleRecord
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("id IN ?", keys).Find(&recs).Error; err != nil {
		return nil, err
	}

	views := make([]query.RoleView, len(recs))
	for i, rec := range recs {
		views[i] = toRoleView(rec.toDomain())
	}
	return views, nil
}

func toRoleView(r *role.Role) query.RoleView {
	perms := r.Permissions()
	names := make([]string, len(perms))
	for i, p := range perms {
		names[i] = p.String()
	}
	return query.RoleView{
		ID:             r.ID().String(),
		OrganizationID: r.OrganizationID().String(),
		Name:           r.Name(),
		Owner:          r.IsOwner(),
		Permissions:    names,
		CreatedAt:      r.CreatedAt(),
	}
}
