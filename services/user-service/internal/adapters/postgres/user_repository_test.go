package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	pgshared "github.com/ChristianPenev/chrisos/pkg/postgres"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/postgres/migrations"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	if err := pgshared.Migrate(context.Background(), sqlDB, migrations.FS); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("TRUNCATE users, role_permissions, roles, organizations").Error; err != nil {
		t.Fatal(err)
	}

	return db
}

func seedOrganization(t *testing.T, db *gorm.DB) (organization.ID, role.ID) {
	t.Helper()
	org, err := organization.NewOrganization("Acme")
	if err != nil {
		t.Fatal(err)
	}
	owner := role.NewOwnerRole(org.ID())

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := NewOrganizationRepository(tx).Create(context.Background(), org); err != nil {
			return err
		}
		return tx.Exec("INSERT INTO roles (id, organization_id, name, owner, created_at) VALUES (?, ?, ?, ?, ?)",
			owner.ID().String(), org.ID().String(), owner.Name(), owner.IsOwner(), owner.CreatedAt()).Error
	})
	if err != nil {
		t.Fatal(err)
	}
	return org.ID(), owner.ID()
}

func newUser(t *testing.T, orgID organization.ID, roleID role.ID, rawEmail string) *user.User {
	t.Helper()
	email, err := user.NewEmail(rawEmail)
	if err != nil {
		t.Fatal(err)
	}
	u, err := user.NewUser(orgID, roleID, email, "Jane", "Doe")
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func assertSameUser(t *testing.T, got, want *user.User) {
	t.Helper()
	if got.ID() != want.ID() ||
		got.OrganizationID() != want.OrganizationID() ||
		got.RoleID() != want.RoleID() ||
		got.Email() != want.Email() ||
		got.FirstName() != want.FirstName() ||
		got.LastName() != want.LastName() ||
		!got.CreatedAt().Equal(want.CreatedAt()) {
		t.Fatalf("round trip mismatch:\n got  %+v\n want %+v", *got, *want)
	}
}

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	repo := NewUserRepository(db)
	orgID, roleID := seedOrganization(t, db)
	jane := newUser(t, orgID, roleID, "jane@example.com")

	if err := repo.Create(ctx, jane); err != nil {
		t.Fatalf("Create: %v", err)
	}

	t.Run("FindByID returns the same user", func(t *testing.T) {
		got, err := repo.FindByID(ctx, jane.ID())
		if err != nil {
			t.Fatal(err)
		}
		assertSameUser(t, got, jane)
	})

	t.Run("FindByEmail returns the same user", func(t *testing.T) {
		got, err := repo.FindByEmail(ctx, jane.Email())
		if err != nil {
			t.Fatal(err)
		}
		assertSameUser(t, got, jane)
	})

	t.Run("Create with taken email", func(t *testing.T) {
		err := repo.Create(ctx, newUser(t, orgID, roleID, "jane@example.com"))
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Fatalf("want ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("FindByID unknown", func(t *testing.T) {
		_, err := repo.FindByID(ctx, user.NewID())
		if !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})

	t.Run("FindByEmail unknown", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, newUser(t, orgID, roleID, "nobody@example.com").Email())
		if !errors.Is(err, user.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})
}
