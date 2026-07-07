package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/infrastructure/auth"
)

const secret = "a-test-secret-of-at-least-32-characters"

func setupAuth(t *testing.T) (*gin.Engine, *auth.JWTManager) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	jwt := auth.NewJWTManager(secret, time.Hour)

	r := gin.New()
	r.GET("/protected", AuthMiddleware(jwt), func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		role, _ := c.Get("user_role")
		c.JSON(http.StatusOK, gin.H{"uid": uid, "role": role})
	})
	r.GET("/admin", AuthMiddleware(jwt), RequireRole("admin"), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r, jwt
}

func TestAuthMiddleware_RejectsMissingHeader(t *testing.T) {
	r, _ := setupAuth(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_RejectsMalformedHeader(t *testing.T) {
	r, _ := setupAuth(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic abc")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	r, _ := setupAuth(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_AcceptsValidToken(t *testing.T) {
	r, jwt := setupAuth(t)
	token, _ := jwt.Generate("user-1", "broadcaster")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestRequireRole_BlocksWrongRole(t *testing.T) {
	r, jwt := setupAuth(t)
	token, _ := jwt.Generate("user-1", "user")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequireRole_AllowsMatchingRole(t *testing.T) {
	r, jwt := setupAuth(t)
	token, _ := jwt.Generate("admin-1", "admin")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
