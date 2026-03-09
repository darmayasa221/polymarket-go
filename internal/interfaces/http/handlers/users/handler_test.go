package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	adduserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser/dto"
	getuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
	listuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
	jsonpresenter "github.com/darmayasa221/polymarket-go/internal/interfaces/http/presenters/json"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// --- mock use cases ---

type mockAddUser struct {
	out adduserDTO.Output
	err error
}

func (m *mockAddUser) Execute(_ context.Context, _ adduserDTO.Input) (adduserDTO.Output, error) {
	return m.out, m.err
}

type mockGetUser struct {
	out getuserDTO.Output
	err error
}

func (m *mockGetUser) Execute(_ context.Context, _ getuserDTO.Input) (getuserDTO.Output, error) {
	return m.out, m.err
}

type mockListUsers struct {
	offsetOut listuserDTO.OffsetOutput
	cursorOut listuserDTO.CursorOutput
	err       error
}

func (m *mockListUsers) ExecuteOffset(_ context.Context, _ listuserDTO.Input) (listuserDTO.OffsetOutput, error) {
	return m.offsetOut, m.err
}

func (m *mockListUsers) ExecuteCursor(_ context.Context, _ listuserDTO.Input) (listuserDTO.CursorOutput, error) {
	return m.cursorOut, m.err
}

// --- helpers ---

func newTestHandler(t *testing.T, addUser *mockAddUser, getUser *mockGetUser, listUsers *mockListUsers) *Handler {
	t.Helper()
	logger, err := logging.New("error")
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	return New(addUser, getUser, listUsers, jsonpresenter.New(), logger)
}

func newGinContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	} else {
		req = httptest.NewRequest(method, path, http.NoBody)
	}
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return m
}

// --- Register tests ---

func TestRegister_InvalidBody(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{})
	c, w := newGinContext(http.MethodPost, "/users", []byte(`{}`))

	h.Register(c)

	// Sending an empty body triggers validator.ValidationErrors → 422 Unprocessable Entity.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestRegister_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{err: errors.New("conflict")}, &mockGetUser{}, &mockListUsers{})
	payload, _ := json.Marshal(map[string]string{
		"username":  "alice",
		"email":     "alice@example.com",
		"password":  "secret",
		"full_name": "Alice",
	})
	c, w := newGinContext(http.MethodPost, "/users", payload)

	h.Register(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRegister_Success(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockAddUser{out: adduserDTO.Output{
		ID: "u1", Username: "alice", Email: "alice@example.com", FullName: "Alice", CreatedAt: now,
	}}, &mockGetUser{}, &mockListUsers{})
	payload, _ := json.Marshal(map[string]string{
		"username":  "alice",
		"email":     "alice@example.com",
		"password":  "secret",
		"full_name": "Alice",
	})
	c, w := newGinContext(http.MethodPost, "/users", payload)

	h.Register(c)

	if w.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusCreated)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}

// --- GetMe tests ---

func TestGetMe_MissingUserID(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{})
	c, w := newGinContext(http.MethodGet, "/users/me", nil)

	h.GetMe(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestGetMe_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{err: errors.New("not found")}, &mockListUsers{})
	c, w := newGinContext(http.MethodGet, "/users/me", nil)
	c.Set(response.ContextKeyUserID, "u1")

	h.GetMe(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetMe_Success(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{out: getuserDTO.Output{
		ID: "u1", Username: "alice", Email: "alice@example.com", FullName: "Alice",
		CreatedAt: now, UpdatedAt: now,
	}}, &mockListUsers{})
	c, w := newGinContext(http.MethodGet, "/users/me", nil)
	c.Set(response.ContextKeyUserID, "u1")

	h.GetMe(c)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}

// --- List tests ---

// --- GetByID tests ---

func TestGetByID_InvalidURI(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{})
	// No URI param set — ShouldBindUri will fail because "id" is required.
	c, w := newGinContext(http.MethodGet, "/users/", nil)

	h.GetByID(c)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestList_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{err: errors.New("db error")})
	c, w := newGinContext(http.MethodGet, "/users", nil)

	h.List(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestList_OffsetSuccess(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{
		offsetOut: listuserDTO.OffsetOutput{
			Users:      []listuserDTO.UserItem{{ID: "u1", Username: "alice", Email: "alice@example.com", FullName: "Alice", CreatedAt: now}},
			Page:       1,
			PageSize:   10,
			TotalItems: 1,
			TotalPages: 1,
		},
	})
	c, w := newGinContext(http.MethodGet, "/users", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}

func TestList_CursorSuccess(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockAddUser{}, &mockGetUser{}, &mockListUsers{
		cursorOut: listuserDTO.CursorOutput{
			Users:      []listuserDTO.UserItem{{ID: "u2", Username: "bob", Email: "bob@example.com", FullName: "Bob", CreatedAt: now}},
			NextCursor: "cursor_next",
			HasNext:    true,
		},
	})
	c, w := newGinContext(http.MethodGet, "/users?cursor=cursor_abc", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}
