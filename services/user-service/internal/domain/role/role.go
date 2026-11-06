package role

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
)

const (
	maxNameLength = 100
	OwnerRoleName = "Organization Owner"
)

type Role struct {
	id             ID
	organizationID organization.ID
	name           string
	permissions    []Permission
	owner          bool
	createdAt      time.Time
}

func NewRole(organizationID organization.ID, name string, permissions []Permission) (*Role, error) {
	n, err := normalizeName(name)
	if err != nil {
		return nil, err
	}
	perms, err := normalizePermissions(permissions)
	if err != nil {
		return nil, err
	}

	return &Role{
		id:             NewID(),
		organizationID: organizationID,
		name:           n,
		permissions:    perms,
		createdAt:      now(),
	}, nil
}

func NewOwnerRole(organizationID organization.ID) *Role {
	return &Role{
		id:             NewID(),
		organizationID: organizationID,
		name:           OwnerRoleName,
		owner:          true,
		createdAt:      now(),
	}
}

func Rehydrate(id ID, organizationID organization.ID, name string, permissions []string, owner bool, createdAt time.Time) *Role {
	perms := make([]Permission, 0, len(permissions))
	for _, s := range permissions {
		if p, err := ParsePermission(s); err == nil {
			perms = append(perms, p)
		}
	}
	slices.Sort(perms)

	return &Role{
		id:             id,
		organizationID: organizationID,
		name:           name,
		permissions:    slices.Compact(perms),
		owner:          owner,
		createdAt:      createdAt.UTC(),
	}
}

func (r *Role) ID() ID                          { return r.id }
func (r *Role) OrganizationID() organization.ID { return r.organizationID }
func (r *Role) Name() string                    { return r.name }
func (r *Role) IsOwner() bool                   { return r.owner }
func (r *Role) CreatedAt() time.Time            { return r.createdAt }

func (r *Role) Permissions() []Permission {
	if r.owner {
		return AllPermissions()
	}
	return slices.Clone(r.permissions)
}

func (r *Role) Has(p Permission) bool {
	return r.owner || slices.Contains(r.permissions, p)
}

func (r *Role) Rename(name string) error {
	if r.owner {
		return ErrOwnerRoleLocked
	}
	n, err := normalizeName(name)
	if err != nil {
		return err
	}
	r.name = n
	return nil
}

func (r *Role) SetPermissions(permissions []Permission) error {
	if r.owner {
		return ErrOwnerRoleLocked
	}
	perms, err := normalizePermissions(permissions)
	if err != nil {
		return err
	}
	r.permissions = perms
	return nil
}

func normalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", fmt.Errorf("%w: must not be empty", ErrInvalidName)
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return "", fmt.Errorf("%w: longer than %d characters", ErrInvalidName, maxNameLength)
	}
	return name, nil
}

func normalizePermissions(permissions []Permission) ([]Permission, error) {
	perms := make([]Permission, 0, len(permissions))
	for _, p := range permissions {
		if _, err := ParsePermission(string(p)); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	slices.Sort(perms)
	return slices.Compact(perms), nil
}

func now() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}
