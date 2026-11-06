package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ChristianPenev/chrisos/apps/people/api/internal/app"
)

type fakeUserService struct {
	users     map[string]*app.User
	createErr error
	getErr    error
	listErr   error
	gotList   app.ListUsersInput
	page      *app.UserPage
}

func (f *fakeUserService) ListUsers(_ context.Context, in app.ListUsersInput) (*app.UserPage, error) {
	f.gotList = in
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.page, nil
}

func (f *fakeUserService) CreateOrganization(_ context.Context, in app.CreateOrganizationInput) (app.CreatedOrganization, error) {
	if f.createErr != nil {
		return app.CreatedOrganization{}, f.createErr
	}
	f.users["user-1"] = &app.User{
		ID: "user-1", OrganizationID: "org-1", RoleID: "role-1",
		Email: in.OwnerEmail, FirstName: in.OwnerFirstName, LastName: in.OwnerLastName,
		CreatedAt: time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC),
	}
	return app.CreatedOrganization{OrganizationID: "org-1", OwnerUserID: "user-1"}, nil
}

func (f *fakeUserService) GetUser(_ context.Context, id string) (*app.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	u, ok := f.users[id]
	if !ok {
		return nil, app.ErrUserNotFound
	}
	return u, nil
}

type gqlResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message    string         `json:"message"`
		Path       []any          `json:"path"`
		Extensions map[string]any `json:"extensions"`
	} `json:"errors"`
}

func do(t *testing.T, fake *fakeUserService, query string) gqlResponse {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"query": query})
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	NewHandler(app.New(fake), true).ServeHTTP(rec, req)

	var resp gqlResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return resp
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

const createAcme = `mutation {
	createOrganization(input: {name: "Acme", ownerEmail: "jane@example.com", ownerFirstName: "Jane", ownerLastName: "Doe"}) {
		organizationId
		owner { id organizationId roleId email firstName lastName createdAt }
	}
}`

func TestCreateOrganization(t *testing.T) {
	resp := do(t, &fakeUserService{users: map[string]*app.User{}}, createAcme)

	if len(resp.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", resp.Errors)
	}
	want := `{"organizationId":"org-1","owner":{"id":"user-1","organizationId":"org-1","roleId":"role-1","email":"jane@example.com","firstName":"Jane","lastName":"Doe","createdAt":"2026-11-06T10:00:00Z"}}`
	if got := string(resp.Data["createOrganization"]); got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestUser(t *testing.T) {
	fake := &fakeUserService{users: map[string]*app.User{
		"user-1": {ID: "user-1", Email: "jane@example.com", FirstName: "Jane", LastName: "Doe"},
	}}
	resp := do(t, fake, `{ user(id: "user-1") { id email } }`)

	if len(resp.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", resp.Errors)
	}
	if got := string(resp.Data["user"]); got != `{"id":"user-1","email":"jane@example.com"}` {
		t.Fatalf("got %s", got)
	}
}

func TestUserNotFound(t *testing.T) {
	resp := do(t, &fakeUserService{users: map[string]*app.User{}}, `{ user(id: "missing") { id } }`)

	assertErrorCode(t, resp, codeNotFound)
	if string(resp.Data["user"]) != "null" {
		t.Errorf("want data.user null, got %s", resp.Data["user"])
	}
	if fmt.Sprint(resp.Errors[0].Path) != "[user]" {
		t.Errorf("want path [user], got %v", resp.Errors[0].Path)
	}
}

func TestCreateOrganizationErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode string
		wantMsg  string
	}{
		{"invalid input", fmt.Errorf("%w: %s", app.ErrInvalidInput, `invalid email: "nope"`), codeBadUserInput, `invalid input: invalid email: "nope"`},
		{"email taken", app.ErrEmailAlreadyExists, codeConflict, "email already exists"},
		{"unexpected", errors.New("user-service: connection refused 10.0.0.7:9090"), codeInternal, "internal error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := do(t, &fakeUserService{users: map[string]*app.User{}, createErr: tt.err}, createAcme)
			assertErrorCode(t, resp, tt.wantCode)
			if msg := resp.Errors[0].Message; msg != tt.wantMsg {
				t.Errorf("want message %q, got %q", tt.wantMsg, msg)
			}
		})
	}
}

func TestUsers(t *testing.T) {
	jane := &app.User{ID: "user-1", Email: "jane@example.com"}
	john := &app.User{ID: "user-2", Email: "john@example.com"}

	t.Run("middle page", func(t *testing.T) {
		fake := &fakeUserService{page: &app.UserPage{Users: []*app.User{jane, john}, NextCursor: "cursor-2"}}
		resp := do(t, fake, `{ users(organizationId: "org-1", first: 2, after: "cursor-1") {
			nodes { id email }
			pageInfo { hasNextPage endCursor }
		} }`)

		if len(resp.Errors) != 0 {
			t.Fatalf("unexpected errors: %+v", resp.Errors)
		}
		want := `{"nodes":[{"id":"user-1","email":"jane@example.com"},{"id":"user-2","email":"john@example.com"}],"pageInfo":{"hasNextPage":true,"endCursor":"cursor-2"}}`
		if got := string(resp.Data["users"]); got != want {
			t.Fatalf("got  %s\nwant %s", got, want)
		}
		if fake.gotList != (app.ListUsersInput{OrganizationID: "org-1", First: 2, After: "cursor-1"}) {
			t.Errorf("service got %+v", fake.gotList)
		}
	})

	t.Run("last page has no end cursor", func(t *testing.T) {
		fake := &fakeUserService{page: &app.UserPage{Users: []*app.User{jane}}}
		resp := do(t, fake, `{ users(organizationId: "org-1") { nodes { id } pageInfo { hasNextPage endCursor } } }`)

		if len(resp.Errors) != 0 {
			t.Fatalf("unexpected errors: %+v", resp.Errors)
		}
		want := `{"nodes":[{"id":"user-1"}],"pageInfo":{"hasNextPage":false,"endCursor":null}}`
		if got := string(resp.Data["users"]); got != want {
			t.Fatalf("got  %s\nwant %s", got, want)
		}
		if fake.gotList != (app.ListUsersInput{OrganizationID: "org-1"}) {
			t.Errorf("omitted arguments should be zero values, service got %+v", fake.gotList)
		}
	})

	t.Run("empty organization", func(t *testing.T) {
		fake := &fakeUserService{page: &app.UserPage{Users: []*app.User{}}}
		resp := do(t, fake, `{ users(organizationId: "org-1") { nodes { id } pageInfo { hasNextPage } } }`)

		if got := string(resp.Data["users"]); got != `{"nodes":[],"pageInfo":{"hasNextPage":false}}` {
			t.Fatalf("got %s", got)
		}
	})

	t.Run("invalid input is BAD_USER_INPUT", func(t *testing.T) {
		fake := &fakeUserService{listErr: fmt.Errorf("%w: %s", app.ErrInvalidInput, "invalid cursor")}
		resp := do(t, fake, `{ users(organizationId: "org-1", after: "garbage") { nodes { id } } }`)
		assertErrorCode(t, resp, codeBadUserInput)
	})
}

func TestValidationErrorsKeepTheirCode(t *testing.T) {
	resp := do(t, &fakeUserService{users: map[string]*app.User{}}, `{ user(id: "x") { password } }`)
	assertErrorCode(t, resp, "GRAPHQL_VALIDATION_FAILED")
}

func TestComplexityLimit(t *testing.T) {
	fields := strings.Repeat("u%d: user(id: \"x\") { id } ", maxQueryComplexity)
	aliases := make([]any, maxQueryComplexity)
	for i := range aliases {
		aliases[i] = i
	}
	resp := do(t, &fakeUserService{users: map[string]*app.User{}}, "{ "+fmt.Sprintf(fields, aliases...)+" }")
	assertErrorCode(t, resp, "COMPLEXITY_LIMIT_EXCEEDED")
}

func TestIntrospectionDisabled(t *testing.T) {
	body := strings.NewReader(`{"query":"{ __schema { queryType { name } } }"}`)
	req := httptest.NewRequest(http.MethodPost, "/graphql", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	NewHandler(app.New(&fakeUserService{users: map[string]*app.User{}}), false).ServeHTTP(rec, req)

	var resp gqlResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Message != "introspection disabled" {
		t.Fatalf("want an introspection disabled error, got %s", rec.Body.String())
	}
}
