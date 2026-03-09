package authentications

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

	loginuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser/dto"
	logoutuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser/dto"
	refreshauthDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
	jsonpresenter "github.com/darmayasa221/polymarket-go/internal/interfaces/http/presenters/json"
)

// --- mock use cases ---

type mockLoginUser struct {
	out loginuserDTO.Output
	err error
}

func (m *mockLoginUser) Execute(_ context.Context, _ loginuserDTO.Input) (loginuserDTO.Output, error) {
	return m.out, m.err
}

type mockLogoutUser struct {
	out       logoutuserDTO.Output
	err       error
	lastInput logoutuserDTO.Input
}

func (m *mockLogoutUser) Execute(_ context.Context, in logoutuserDTO.Input) (logoutuserDTO.Output, error) {
	m.lastInput = in
	return m.out, m.err
}

type mockRefreshAuth struct {
	out refreshauthDTO.Output
	err error
}

func (m *mockRefreshAuth) Execute(_ context.Context, _ refreshauthDTO.Input) (refreshauthDTO.Output, error) {
	return m.out, m.err
}

// --- helpers ---

func newTestHandler(t *testing.T, loginUser *mockLoginUser, logoutUser *mockLogoutUser, refreshAuth *mockRefreshAuth) *Handler {
	t.Helper()
	logger, err := logging.New("error")
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	return New(loginUser, logoutUser, refreshAuth, jsonpresenter.New(), logger)
}

func newGinContext(path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
		req.Header.Set(httpconst.HeaderContentType, httpconst.ContentTypeJSON)
	} else {
		req = httptest.NewRequest(http.MethodPost, path, http.NoBody)
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

// --- Login tests ---

func TestLogin_InvalidBody(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{}, &mockRefreshAuth{})
	c, w := newGinContext("/auth/login", []byte(`{}`))

	h.Login(c)

	// Sending an empty body triggers validator.ValidationErrors → 422 Unprocessable Entity.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestLogin_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{err: errors.New("invalid credentials")}, &mockLogoutUser{}, &mockRefreshAuth{})
	payload, _ := json.Marshal(map[string]string{
		"username": "alice",
		"password": "wrongpass",
	})
	c, w := newGinContext("/auth/login", payload)

	h.Login(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestLogin_Success(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockLoginUser{out: loginuserDTO.Output{
		AccessToken:           "access-tok",
		RefreshToken:          "refresh-tok",
		AccessTokenExpiresAt:  now,
		RefreshTokenExpiresAt: now,
		UserID:                "u1",
	}}, &mockLogoutUser{}, &mockRefreshAuth{})
	payload, _ := json.Marshal(map[string]string{
		"username": "alice",
		"password": "secret",
	})
	c, w := newGinContext("/auth/login", payload)

	h.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}

// --- Logout tests ---

func TestLogout_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{err: errors.New("token not found")}, &mockRefreshAuth{})
	c, w := newGinContext("/auth/logout", nil)
	c.Request.Header.Set(httpconst.HeaderAuthorization, httpconst.PrefixBearer+"some-token")

	h.Logout(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestLogout_Success(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{}, &mockRefreshAuth{})
	c, w := newGinContext("/auth/logout", nil)
	c.Request.Header.Set(httpconst.HeaderAuthorization, httpconst.PrefixBearer+"some-token")

	h.Logout(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestLogout_MissingAuthHeader(t *testing.T) {
	mock := &mockLogoutUser{}
	h := newTestHandler(t, &mockLoginUser{}, mock, &mockRefreshAuth{})
	c, w := newGinContext("/auth/logout", nil)

	// No Authorization header — logout proceeds with empty TokenValue (idempotent).
	h.Logout(c)
	c.Writer.WriteHeaderNow()

	if w.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusNoContent)
	}
	if mock.lastInput.TokenValue != "" {
		t.Errorf("TokenValue: got %q, want empty string", mock.lastInput.TokenValue)
	}
}

// --- Refresh tests ---

func TestRefresh_InvalidBody(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{}, &mockRefreshAuth{})
	c, w := newGinContext("/auth/refresh", []byte(`{}`))

	h.Refresh(c)

	// Sending an empty body triggers validator.ValidationErrors → 422 Unprocessable Entity.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestRefresh_UseCaseError(t *testing.T) {
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{}, &mockRefreshAuth{err: errors.New("token expired")})
	payload, _ := json.Marshal(map[string]string{
		"refresh_token": "expired-tok",
	})
	c, w := newGinContext("/auth/refresh", payload)

	h.Refresh(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != false {
		t.Errorf("success: got %v, want false", body["success"])
	}
}

func TestRefresh_Success(t *testing.T) {
	now := time.Now()
	h := newTestHandler(t, &mockLoginUser{}, &mockLogoutUser{}, &mockRefreshAuth{out: refreshauthDTO.Output{
		AccessToken:           "new-access-tok",
		RefreshToken:          "new-refresh-tok",
		AccessTokenExpiresAt:  now,
		RefreshTokenExpiresAt: now,
	}})
	payload, _ := json.Marshal(map[string]string{
		"refresh_token": "valid-refresh-tok",
	})
	c, w := newGinContext("/auth/refresh", payload)

	h.Refresh(c)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeBody(t, w.Body.Bytes())
	if body["success"] != true {
		t.Errorf("success: got %v, want true", body["success"])
	}
}
