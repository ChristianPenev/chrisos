package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/ChristianPenev/chrisos/services/user-service/internal/adapters/memory"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/command"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/application/query"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/organization"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/role"
	"github.com/ChristianPenev/chrisos/services/user-service/internal/domain/user"
)

func newApp(store *memory.Store) application.Application {
	return application.Application{
		Commands: application.Commands{CreateOrganization: command.NewCreateOrganizationHandler(store)},
		Queries: application.Queries{
			GetUser:   query.NewGetUserHandler(store),
			GetUsers:  query.NewGetUsersHandler(store),
			ListUsers: query.NewListUsersHandler(store),
			GetRoles:  query.NewGetRolesHandler(store),
		},
	}
}

type gqlResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message    string         `json:"message"`
		Path       []any          `json:"path"`
		Extensions map[string]any `json:"extensions"`
	} `json:"errors"`
}

func doWith(t *testing.T, store *memory.Store, introspection bool, query string, variables map[string]any) gqlResponse {
	t.Helper()

	body, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	NewHandler(newApp(store), introspection).ServeHTTP(rec, req)

	var resp gqlResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return resp
}

func do(t *testing.T, store *memory.Store, query string, variables map[string]any) gqlResponse {
	t.Helper()
	return doWith(t, store, true, query, variables)
}

func assertNoErrors(t *testing.T, resp gqlResponse) {
	t.Helper()
	if len(resp.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", resp.Errors)
	}
}

func assertErrorCode(t *testing.T, resp gqlResponse, want string) {
	t.Helper()
	if len(resp.Errors) != 1 {
		t.Fatalf("want 1 error, got %+v", resp.Errors)
	}
	if got := resp.Errors[0].Extensions["code"]; got != want {
		t.Fatalf("want code %s, got %v (%s)", want, got, resp.Errors[0].Message)
	}
}

func decode[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("decode %s: %v", raw, err)
	}
	return v
}

func seedOrganization(t *testing.T, store *memory.Store) (organization.ID, *role.Role) {
	t.Helper()
	ctx := context.Background()

	org, err := organization.NewOrganization("Acme")
	if err != nil {
		t.Fatal(err)
	}
	owner := role.NewOwnerRole(org.ID())
	repos := store.Repositories()
	if err := repos.Organizations.Create(ctx, org); err != nil {
		t.Fatal(err)
	}
	if err := repos.Roles.Create(ctx, owner); err != nil {
		t.Fatal(err)
	}
	return org.ID(), owner
}

func seedRole(t *testing.T, store *memory.Store, orgID organization.ID, name string) *role.Role {
	t.Helper()
	r, err := role.NewRole(orgID, name, []role.Permission{role.PermissionMembersInvite})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Repositories().Roles.Create(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	return r
}

func seedUser(t *testing.T, store *memory.Store, orgID organization.ID, roleID role.ID, email string) *user.User {
	t.Helper()
	e, err := user.NewEmail(email)
	if err != nil {
		t.Fatal(err)
	}
	u, err := user.NewUser(orgID, roleID, e, "First", "Last")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Repositories().Users.Create(context.Background(), u); err != nil {
		t.Fatal(err)
	}
	return u
}

const createAcme = `mutation($email: String!) {
	createOrganization(input: {name: "Acme", ownerEmail: $email, ownerFirstName: "Jane", ownerLastName: "Doe"}) {
		organizationId
		owner { id organizationId email firstName lastName createdAt role { name isOwner permissions } }
	}
}`

func TestCreateOrganization(t *testing.T) {
	resp := do(t, memory.NewStore(), createAcme, map[string]any{"email": "Jane@Example.com"})
	assertNoErrors(t, resp)

	got := decode[struct {
		OrganizationID string
		Owner          struct {
			ID, OrganizationID, Email, FirstName, LastName, CreatedAt string
			Role                                                      struct {
				Name        string
				IsOwner     bool
				Permissions []string
			}
		}
	}](t, resp.Data["createOrganization"])

	if got.OrganizationID == "" || got.Owner.OrganizationID != got.OrganizationID {
		t.Fatalf("organization IDs: %+v", got)
	}
	if got.Owner.ID == "" || got.Owner.Email != "jane@example.com" ||
		got.Owner.FirstName != "Jane" || got.Owner.LastName != "Doe" || got.Owner.CreatedAt == "" {
		t.Fatalf("owner: %+v", got.Owner)
	}
	if got.Owner.Role.Name != role.OwnerRoleName || !got.Owner.Role.IsOwner ||
		len(got.Owner.Role.Permissions) != len(role.AllPermissions()) {
		t.Fatalf("owner role: %+v", got.Owner.Role)
	}
}

func TestCreateOrganizationErrors(t *testing.T) {
	t.Run("invalid email is BAD_USER_INPUT", func(t *testing.T) {
		resp := do(t, memory.NewStore(), createAcme, map[string]any{"email": "nope"})
		assertErrorCode(t, resp, codeBadUserInput)
	})

	t.Run("taken email is CONFLICT", func(t *testing.T) {
		store := memory.NewStore()
		assertNoErrors(t, do(t, store, createAcme, map[string]any{"email": "jane@example.com"}))
		resp := do(t, store, createAcme, map[string]any{"email": "jane@example.com"})
		assertErrorCode(t, resp, codeConflict)
	})

	t.Run("unexpected failure is INTERNAL without details", func(t *testing.T) {
		store := memory.NewStore()
		store.UserCreateErr = errors.New("connection refused to db-host:5432")
		resp := do(t, store, createAcme, map[string]any{"email": "jane@example.com"})
		assertErrorCode(t, resp, codeInternal)
		if resp.Errors[0].Message != "internal error" {
			t.Errorf("internal details leaked: %q", resp.Errors[0].Message)
		}
	})
}

func TestUser(t *testing.T) {
	store := memory.NewStore()
	orgID, owner := seedOrganization(t, store)
	jane := seedUser(t, store, orgID, owner.ID(), "jane@example.com")

	t.Run("found", func(t *testing.T) {
		resp := do(t, store, `query($id: ID!) { user(id: $id) { id email roleId } }`, map[string]any{"id": jane.ID().String()})
		assertNoErrors(t, resp)
		want := fmt.Sprintf(`{"id":%q,"email":"jane@example.com","roleId":%q}`, jane.ID(), owner.ID())
		if got := string(resp.Data["user"]); got != want {
			t.Fatalf("got  %s\nwant %s", got, want)
		}
	})

	t.Run("unknown ID is NOT_FOUND", func(t *testing.T) {
		resp := do(t, store, `query($id: ID!) { user(id: $id) { id } }`, map[string]any{"id": user.NewID().String()})
		assertErrorCode(t, resp, codeNotFound)
		if string(resp.Data["user"]) != "null" {
			t.Errorf("want data.user null, got %s", resp.Data["user"])
		}
		if fmt.Sprint(resp.Errors[0].Path) != "[user]" {
			t.Errorf("want path [user], got %v", resp.Errors[0].Path)
		}
	})

	t.Run("malformed ID is BAD_USER_INPUT", func(t *testing.T) {
		resp := do(t, store, `{ user(id: "nope") { id } }`, nil)
		assertErrorCode(t, resp, codeBadUserInput)
	})
}

type userConnection struct {
	Nodes    []struct{ ID string }
	PageInfo struct {
		HasNextPage bool
		EndCursor   *string
	}
}

func TestUsers(t *testing.T) {
	store := memory.NewStore()
	orgID, owner := seedOrganization(t, store)
	var want []string
	for i := range 5 {
		want = append(want, seedUser(t, store, orgID, owner.ID(), fmt.Sprintf("user%d@acme.com", i)).ID().String())
	}
	otherOrg, otherOwner := seedOrganization(t, store)
	seedUser(t, store, otherOrg, otherOwner.ID(), "someone@other.com")
	slices.Sort(want)

	const usersQuery = `query($org: ID!, $first: Int, $after: String) {
		users(organizationId: $org, first: $first, after: $after) {
			nodes { id }
			pageInfo { hasNextPage endCursor }
		}
	}`

	t.Run("pages through the organization's users", func(t *testing.T) {
		var got []string
		var after any
		for range 10 {
			resp := do(t, store, usersQuery, map[string]any{"org": orgID.String(), "first": 2, "after": after})
			assertNoErrors(t, resp)
			page := decode[userConnection](t, resp.Data["users"])
			for _, n := range page.Nodes {
				got = append(got, n.ID)
			}
			if !page.PageInfo.HasNextPage {
				if page.PageInfo.EndCursor != nil {
					t.Errorf("last page endCursor = %q, want null", *page.PageInfo.EndCursor)
				}
				break
			}
			after = *page.PageInfo.EndCursor
		}
		if !slices.Equal(got, want) {
			t.Fatalf("got  %v\nwant %v", got, want)
		}
	})

	t.Run("empty organization", func(t *testing.T) {
		resp := do(t, store, usersQuery, map[string]any{"org": organization.NewID().String()})
		assertNoErrors(t, resp)
		if got := string(resp.Data["users"]); got != `{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}` {
			t.Fatalf("got %s", got)
		}
	})

	invalid := []struct {
		name string
		vars map[string]any
	}{
		{"malformed organization ID", map[string]any{"org": "nope"}},
		{"negative first", map[string]any{"org": orgID.String(), "first": -1}},
		{"garbage cursor", map[string]any{"org": orgID.String(), "after": "garbage"}},
	}
	for _, tt := range invalid {
		t.Run(tt.name+" is BAD_USER_INPUT", func(t *testing.T) {
			assertErrorCode(t, do(t, store, usersQuery, tt.vars), codeBadUserInput)
		})
	}
}

func TestUserRoles(t *testing.T) {
	store := memory.NewStore()
	orgID, owner := seedOrganization(t, store)
	support := seedRole(t, store, orgID, "Support")

	wantRole := map[string]string{}
	for i := range 20 {
		r := support
		if i%5 == 0 {
			r = owner
		}
		u := seedUser(t, store, orgID, r.ID(), fmt.Sprintf("user%d@acme.com", i))
		wantRole[u.ID().String()] = r.ID().String()
	}

	t.Run("one batched role lookup for a whole list", func(t *testing.T) {
		before := len(store.RoleBatches())
		resp := do(t, store, `query($org: ID!) { users(organizationId: $org, first: 20) { nodes { id role { id name isOwner } } } }`,
			map[string]any{"org": orgID.String()})
		assertNoErrors(t, resp)

		batches := store.RoleBatches()[before:]
		if len(batches) != 1 {
			t.Fatalf("want 1 RolesByIDs call, got %d: %v", len(batches), batches)
		}
		var got []string
		for _, id := range batches[0] {
			got = append(got, id.String())
		}
		slices.Sort(got)
		if want := slices.Sorted(slices.Values([]string{owner.ID().String(), support.ID().String()})); !slices.Equal(got, want) {
			t.Fatalf("RolesByIDs got %v, want each role once", batches[0])
		}

		page := decode[struct {
			Nodes []struct {
				ID   string
				Role struct {
					ID      string
					Name    string
					IsOwner bool
				}
			}
		}](t, resp.Data["users"])
		if len(page.Nodes) != 20 {
			t.Fatalf("got %d users", len(page.Nodes))
		}
		for _, n := range page.Nodes {
			if n.Role.ID != wantRole[n.ID] {
				t.Fatalf("user %s got role %s, want %s", n.ID, n.Role.ID, wantRole[n.ID])
			}
			if n.Role.IsOwner != (n.Role.ID == owner.ID().String()) {
				t.Fatalf("role %+v mapped wrong", n.Role)
			}
		}
	})

	t.Run("each request gets fresh loaders", func(t *testing.T) {
		jane := seedUser(t, store, orgID, owner.ID(), "jane@acme.com")
		before := len(store.RoleBatches())
		vars := map[string]any{"id": jane.ID().String()}
		do(t, store, `query($id: ID!) { user(id: $id) { role { name } } }`, vars)
		do(t, store, `query($id: ID!) { user(id: $id) { role { name } } }`, vars)

		if got := len(store.RoleBatches()) - before; got != 2 {
			t.Fatalf("want 2 RolesByIDs calls across 2 requests, got %d", got)
		}
	})

	t.Run("missing role is INTERNAL", func(t *testing.T) {
		orphan := seedUser(t, store, orgID, role.NewID(), "orphan@acme.com")
		resp := do(t, store, `query($id: ID!) { user(id: $id) { id role { name } } }`, map[string]any{"id": orphan.ID().String()})
		assertErrorCode(t, resp, codeInternal)
	})

	t.Run("role lookup failure is INTERNAL without details", func(t *testing.T) {
		failing := memory.NewStore()
		org, ownerRole := seedOrganization(t, failing)
		u := seedUser(t, failing, org, ownerRole.ID(), "jane@acme.com")
		failing.RolesByIDsErr = errors.New("connection refused to db-host:5432")

		resp := do(t, failing, `query($id: ID!) { user(id: $id) { role { name } } }`, map[string]any{"id": u.ID().String()})
		assertErrorCode(t, resp, codeInternal)
		if resp.Errors[0].Message != "internal error" {
			t.Errorf("internal details leaked: %q", resp.Errors[0].Message)
		}
	})
}

const entitiesQuery = `query($reps: [_Any!]!) {
	_entities(representations: $reps) {
		__typename
		... on User { id email role { name } }
		... on Role { id name }
	}
}`

func rep(typename, id string) map[string]any {
	return map[string]any{"__typename": typename, "id": id}
}

func TestEntities(t *testing.T) {
	store := memory.NewStore()
	orgID, owner := seedOrganization(t, store)
	jane := seedUser(t, store, orgID, owner.ID(), "jane@acme.com")
	john := seedUser(t, store, orgID, owner.ID(), "john@acme.com")

	t.Run("resolves mixed representations in order, unknown ones as null", func(t *testing.T) {
		before := len(store.UserBatches())
		resp := do(t, store, entitiesQuery, map[string]any{"reps": []any{
			rep("User", jane.ID().String()),
			rep("Role", owner.ID().String()),
			rep("User", user.NewID().String()),
			rep("User", john.ID().String()),
		}})
		assertNoErrors(t, resp)

		got := decode[[]*struct {
			Typename string `json:"__typename"`
			ID       string
			Email    string
			Name     string
			Role     *struct{ Name string }
		}](t, resp.Data["_entities"])
		if len(got) != 4 {
			t.Fatalf("got %d entities, want 4", len(got))
		}
		if got[0] == nil || got[0].ID != jane.ID().String() || got[0].Email != "jane@acme.com" || got[0].Role.Name != role.OwnerRoleName {
			t.Errorf("entity 0 = %+v, want jane", got[0])
		}
		if got[1] == nil || got[1].Typename != "Role" || got[1].ID != owner.ID().String() || got[1].Name != role.OwnerRoleName {
			t.Errorf("entity 1 = %+v, want the owner role", got[1])
		}
		if got[2] != nil {
			t.Errorf("entity 2 = %+v, want null", got[2])
		}
		if got[3] == nil || got[3].ID != john.ID().String() {
			t.Errorf("entity 3 = %+v, want john", got[3])
		}

		if batches := store.UserBatches()[before:]; len(batches) != 1 || len(batches[0]) != 3 {
			t.Fatalf("want 1 UsersByIDs call with 3 IDs, got %v", batches)
		}
	})

	t.Run("splits more than 100 representations into batches", func(t *testing.T) {
		reps := make([]any, 150)
		for i := range reps {
			reps[i] = rep("User", user.NewID().String())
		}
		before := len(store.UserBatches())
		resp := do(t, store, `query($reps: [_Any!]!) { _entities(representations: $reps) { ... on User { id } } }`,
			map[string]any{"reps": reps})
		assertNoErrors(t, resp)

		var sizes []int
		for _, b := range store.UserBatches()[before:] {
			sizes = append(sizes, len(b))
		}
		if slices.Sort(sizes); !slices.Equal(sizes, []int{50, 100}) {
			t.Fatalf("batch sizes = %v, want [50 100]", sizes)
		}
	})

	t.Run("malformed ID is BAD_USER_INPUT", func(t *testing.T) {
		resp := do(t, store, entitiesQuery, map[string]any{"reps": []any{rep("User", "nope")}})
		assertErrorCode(t, resp, codeBadUserInput)
	})
}

func TestServiceSDL(t *testing.T) {
	resp := do(t, memory.NewStore(), `{ _service { sdl } }`, nil)
	assertNoErrors(t, resp)

	sdl := decode[struct{ SDL string }](t, resp.Data["_service"]).SDL
	for _, want := range []string{
		`specs.apollo.dev/federation/v2`,
		`type User @key(fields: "id")`,
		`type Role @key(fields: "id")`,
		`type PageInfo @shareable`,
	} {
		if !strings.Contains(sdl, want) {
			t.Errorf("sdl is missing %q", want)
		}
	}
}

func TestValidationErrorsKeepTheirCode(t *testing.T) {
	resp := do(t, memory.NewStore(), `{ user(id: "x") { password } }`, nil)
	assertErrorCode(t, resp, "GRAPHQL_VALIDATION_FAILED")
}

func TestComplexityLimit(t *testing.T) {
	var b strings.Builder
	b.WriteString("{ ")
	for i := range maxQueryComplexity {
		fmt.Fprintf(&b, "u%d: user(id: \"x\") { id } ", i)
	}
	b.WriteString("}")
	assertErrorCode(t, do(t, memory.NewStore(), b.String(), nil), "COMPLEXITY_LIMIT_EXCEEDED")
}

func TestIntrospectionDisabled(t *testing.T) {
	resp := doWith(t, memory.NewStore(), false, `{ __schema { queryType { name } } }`, nil)
	if len(resp.Errors) != 1 || resp.Errors[0].Message != "introspection disabled" {
		t.Fatalf("want an introspection disabled error, got %+v", resp.Errors)
	}
}
