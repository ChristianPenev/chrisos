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
	users       map[string]*app.User
	registerErr error
	getErr      error
}

func (f *fakeUserService) RegisterUser(_ context.Context, in app.RegisterUserInput) (string, error) {
	if f.registerErr != nil {
		return "", f.registerErr
	}
	f.users["user-1"] = &app.User{
		ID: "user-1", Email: in.Email, FirstName: in.FirstName, LastName: in.LastName,
		CreatedAt: time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC),
	}
	return "user-1", nil
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

	NewHandler(app.New(fake)).ServeHTTP(rec, req)

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

const registerJane = `mutation { registerUser(input: {email: "jane@example.com", firstName: "Jane", lastName: "Doe"}) { id email firstName lastName createdAt } }`

func TestRegisterUser(t *testing.T) {
	resp := do(t, &fakeUserService{users: map[string]*app.User{}}, registerJane)

	if len(resp.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", resp.Errors)
	}
	want := `{"id":"user-1","email":"jane@example.com","firstName":"Jane","lastName":"Doe","createdAt":"2026-11-06T10:00:00Z"}`
	if got := string(resp.Data["registerUser"]); got != want {
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

func TestRegisterUserErrors(t *testing.T) {
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
			resp := do(t, &fakeUserService{users: map[string]*app.User{}, registerErr: tt.err}, registerJane)
			assertErrorCode(t, resp, tt.wantCode)
			if msg := resp.Errors[0].Message; msg != tt.wantMsg {
				t.Errorf("want message %q, got %q", tt.wantMsg, msg)
			}
		})
	}
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
