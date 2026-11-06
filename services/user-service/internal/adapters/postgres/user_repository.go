package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

const (
	pgUniqueViolation    = "23505"
	usersEmailConstraint = "users_email_key"
)

type UserRepository struct {
	db *gorm.DB
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	rec := userFromDomain(u)
	err := r.db.WithContext(ctx).Create(&rec).Error

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation && pgErr.ConstraintName == usersEmailConstraint {
		return user.ErrEmailAlreadyExists
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id user.ID) (*user.User, error) {
	return r.findOne(ctx, "id = ?", uuid.UUID(id))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	return r.findOne(ctx, "email = ?", email.String())
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg any) (*user.User, error) {
	var rec userRecord
	err := r.db.WithContext(ctx).Where(query, arg).Take(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return rec.toDomain(), nil
}
