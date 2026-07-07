package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/backend/internal/application/dto"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/auth"
)

// memoryUserRepo is a minimal in-memory UserRepository for handler-level tests.
// We deliberately avoid mockgen here to keep the test self-contained.
type memoryUserRepo struct {
	mu       sync.Mutex
	byEmail  map[string]*entity.User
	byID     map[string]*entity.User
	nextID   int
}

func newMemoryUserRepo() *memoryUserRepo {
	return &memoryUserRepo{
		byEmail: map[string]*entity.User{},
		byID:    map[string]*entity.User{},
	}
}

func (r *memoryUserRepo) Create(_ context.Context, u *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEmail[u.Email]; ok {
		return repository.ErrConflict
	}
	r.nextID++
	if u.ID == "" {
		u.ID = "user-" + itoa(r.nextID)
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	r.byEmail[u.Email] = u
	r.byID[u.ID] = u
	return nil
}
func (r *memoryUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}
func (r *memoryUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}
func (r *memoryUserRepo) Update(_ context.Context, u *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}
func (r *memoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[id]; ok {
		delete(r.byEmail, u.Email)
		delete(r.byID, id)
	}
	return nil
}
func (r *memoryUserRepo) List(_ context.Context, offset, limit int) ([]entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]entity.User, 0, len(r.byID))
	for _, u := range r.byID {
		out = append(out, *u)
	}
	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

func setupAuth(t *testing.T) (*gin.Engine, *AuthHandler) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := newMemoryUserRepo()
	jwt := auth.NewJWTManager("a-test-secret-of-at-least-32-characters", time.Hour)
	hasher := auth.NewBcryptHasher()
	uc := usecase.NewAuthUseCase(repo, jwt, hasher)
	h := NewAuthHandler(uc)

	r := gin.New()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	return r, h
}

func TestAuthHandler_RegisterNominal(t *testing.T) {
	r, _ := setupAuth(t)

	body, _ := json.Marshal(dto.RegisterRequest{
		Email:    "alice@example.com",
		Username: "alice",
		Password: "supersecret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%s)", w.Code, w.Body.String())
	}
	var resp dto.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("missing token")
	}
	if resp.User.Email != "alice@example.com" {
		t.Fatalf("unexpected email %q", resp.User.Email)
	}
	if resp.User.Role != string(entity.RoleUser) {
		t.Fatalf("unexpected role %q", resp.User.Role)
	}
}

func TestAuthHandler_RegisterRejectsShortPassword(t *testing.T) {
	r, _ := setupAuth(t)
	body, _ := json.Marshal(dto.RegisterRequest{
		Email:    "bob@example.com",
		Username: "bob",
		Password: "short",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body=%s)", w.Code, w.Body.String())
	}
}

func TestAuthHandler_RegisterRejectsDuplicateEmail(t *testing.T) {
	r, _ := setupAuth(t)
	body, _ := json.Marshal(dto.RegisterRequest{
		Email:    "carol@example.com",
		Username: "carol",
		Password: "supersecret123",
	})
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if i == 0 && w.Code != http.StatusCreated {
			t.Fatalf("first register: expected 201, got %d", w.Code)
		}
		if i == 1 && w.Code != http.StatusConflict {
			t.Fatalf("duplicate register: expected 409, got %d", w.Code)
		}
	}
}

func TestAuthHandler_LoginNominal(t *testing.T) {
	r, _ := setupAuth(t)

	// register first
	reg, _ := json.Marshal(dto.RegisterRequest{
		Email: "dave@example.com", Username: "dave", Password: "supersecret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reg))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	// then login
	body, _ := json.Marshal(dto.LoginRequest{Email: "dave@example.com", Password: "supersecret123"})
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", w.Code, w.Body.String())
	}
	var resp dto.AuthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Fatal("missing token")
	}
}

func TestAuthHandler_LoginRejectsBadPassword(t *testing.T) {
	r, _ := setupAuth(t)
	reg, _ := json.Marshal(dto.RegisterRequest{
		Email: "eve@example.com", Username: "eve", Password: "supersecret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reg))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	body, _ := json.Marshal(dto.LoginRequest{Email: "eve@example.com", Password: "wrong-pass"})
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body=%s)", w.Code, w.Body.String())
	}
}

func TestAuthHandler_LoginRejectsUnknownUser(t *testing.T) {
	r, _ := setupAuth(t)
	body, _ := json.Marshal(dto.LoginRequest{Email: "ghost@example.com", Password: "doesntmatter"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
