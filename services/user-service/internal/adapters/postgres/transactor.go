package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
)

type Transactor struct {
	db *gorm.DB
}

var _ command.Transactor = (*Transactor)(nil)

func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) WithinTransaction(ctx context.Context, fn func(command.Repositories) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(command.Repositories{
			Organizations: NewOrganizationRepository(tx),
			Roles:         NewRoleRepository(tx),
			Users:         NewUserRepository(tx),
		})
	})
}
