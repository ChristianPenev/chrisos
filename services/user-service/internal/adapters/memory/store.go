package memory

import (
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

type Store struct {
	UserCreateErr error
	RolesByIDsErr error

	mu          sync.Mutex
	orgs        map[organization.ID]*organization.Organization
	roles       map[role.ID]*role.Role
	users       map[user.ID]*user.User
	userBatches [][]user.ID
	roleBatches [][]role.ID
}

var (
	_ command.Transactor = (*Store)(nil)
	_ query.UserReader   = (*Store)(nil)
	_ query.RoleReader   = (*Store)(nil)
)

func NewStore() *Store {
	return &Store{
		orgs:  map[organization.ID]*organization.Organization{},
		roles: map[role.ID]*role.Role{},
		users: map[user.ID]*user.User{},
	}
}

func (s *Store) WithinTransaction(_ context.Context, fn func(command.Repositories) error) error {
	return fn(s.Repositories())
}

func (s *Store) Repositories() command.Repositories {
	return command.Repositories{
		Organizations: organizations{s},
		Roles:         roles{s},
		Users:         users{s},
	}
}

func (s *Store) Role(id role.ID) (*role.Role, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.roles[id]
	return r, ok
}

func (s *Store) UserBatches() [][]user.ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.userBatches)
}

func (s *Store) RoleBatches() [][]role.ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.roleBatches)
}

func (s *Store) UserByID(_ context.Context, id user.ID) (query.UserView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return query.UserView{}, user.ErrUserNotFound
	}
	return toUserView(u), nil
}

func (s *Store) UsersByIDs(_ context.Context, ids []user.ID) ([]query.UserView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userBatches = append(s.userBatches, slices.Clone(ids))

	var views []query.UserView
	for _, id := range ids {
		if u, ok := s.users[id]; ok {
			views = append(views, toUserView(u))
		}
	}
	return views, nil
}

func (s *Store) UsersByOrganization(_ context.Context, orgID organization.ID, req pagination.Request) (pagination.Page[query.UserView], error) {
	after := ""
	if req.After != "" {
		key, err := pagination.DecodeCursor[map[string]string](req.After)
		if err != nil {
			return pagination.Page[query.UserView]{}, err
		}
		after = key["id"]
	}

	s.mu.Lock()
	var views []query.UserView
	for _, u := range s.users {
		if u.OrganizationID() == orgID && u.ID().String() > after {
			views = append(views, toUserView(u))
		}
	}
	s.mu.Unlock()
	slices.SortFunc(views, func(a, b query.UserView) int { return strings.Compare(a.ID, b.ID) })

	if len(views) <= req.Size {
		return pagination.Page[query.UserView]{Items: views}, nil
	}
	views = views[:req.Size]
	next, err := pagination.EncodeCursor(map[string]string{"id": views[len(views)-1].ID})
	if err != nil {
		return pagination.Page[query.UserView]{}, err
	}
	return pagination.Page[query.UserView]{Items: views, NextCursor: next}, nil
}

func (s *Store) RolesByIDs(_ context.Context, ids []role.ID) ([]query.RoleView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.roleBatches = append(s.roleBatches, slices.Clone(ids))
	if s.RolesByIDsErr != nil {
		return nil, s.RolesByIDsErr
	}

	var views []query.RoleView
	for _, id := range ids {
		if r, ok := s.roles[id]; ok {
			views = append(views, toRoleView(r))
		}
	}
	return views, nil
}

func toUserView(u *user.User) query.UserView {
	return query.UserView{
		ID:             u.ID().String(),
		OrganizationID: u.OrganizationID().String(),
		RoleID:         u.RoleID().String(),
		Email:          u.Email().String(),
		FirstName:      u.FirstName(),
		LastName:       u.LastName(),
		CreatedAt:      u.CreatedAt(),
	}
}

func toRoleView(r *role.Role) query.RoleView {
	perms := make([]string, 0, len(r.Permissions()))
	for _, p := range r.Permissions() {
		perms = append(perms, p.String())
	}
	return query.RoleView{
		ID:             r.ID().String(),
		OrganizationID: r.OrganizationID().String(),
		Name:           r.Name(),
		Owner:          r.IsOwner(),
		Permissions:    perms,
		CreatedAt:      r.CreatedAt(),
	}
}

type organizations struct{ s *Store }

func (r organizations) Create(_ context.Context, o *organization.Organization) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.orgs[o.ID()] = o
	return nil
}

func (r organizations) FindByID(_ context.Context, id organization.ID) (*organization.Organization, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.orgs[id]
	if !ok {
		return nil, organization.ErrOrganizationNotFound
	}
	return o, nil
}

type roles struct{ s *Store }

func (r roles) Create(_ context.Context, ro *role.Role) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.roles[ro.ID()] = ro
	return nil
}

func (r roles) FindByID(_ context.Context, id role.ID) (*role.Role, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	ro, ok := r.s.roles[id]
	if !ok {
		return nil, role.ErrRoleNotFound
	}
	return ro, nil
}

type users struct{ s *Store }

func (r users) Create(_ context.Context, u *user.User) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if r.s.UserCreateErr != nil {
		return r.s.UserCreateErr
	}
	for _, existing := range r.s.users {
		if existing.Email() == u.Email() {
			return user.ErrEmailAlreadyExists
		}
	}
	r.s.users[u.ID()] = u
	return nil
}

func (r users) FindByID(_ context.Context, id user.ID) (*user.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	u, ok := r.s.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (r users) FindByEmail(_ context.Context, email user.Email) (*user.User, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, u := range r.s.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, user.ErrUserNotFound
}
