package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// memUserRepo is a tiny in-memory UserRepository for handler tests.
type memUserRepo struct {
	byID    map[string]*entity.User
	byEmail map[string]*entity.User
}

func newMemUserRepo() *memUserRepo {
	return &memUserRepo{byID: map[string]*entity.User{}, byEmail: map[string]*entity.User{}}
}

func (r *memUserRepo) Create(_ context.Context, u *entity.User) error {
	if _, exists := r.byEmail[u.Email]; exists {
		return assertDuplicate
	}
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}
func (r *memUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *memUserRepo) FindByEmail(_ context.Context, e string) (*entity.User, error) {
	if u, ok := r.byEmail[e]; ok {
		return u, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *memUserRepo) Update(_ context.Context, u *entity.User) error { r.byID[u.ID] = u; return nil }
func (r *memUserRepo) Delete(_ context.Context, id string) error      { delete(r.byID, id); return nil }
func (r *memUserRepo) List(_ context.Context, _, _ int) ([]entity.User, error) {
	return nil, nil
}

var assertDuplicate = &dupErr{}

type dupErr struct{}

func (*dupErr) Error() string { return "duplicate" }

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	uc := usecase.NewAuthUseCase(newMemUserRepo(), auth.NewJWTManager("secret", time.Hour), auth.NewBcryptHasher())
	h := NewAuthHandler(uc)
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	return r
}

func doJSON(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthHandler_RegisterThenLogin(t *testing.T) {
	r := setupAuthRouter()

	// Register
	w := doJSON(r, http.MethodPost, "/register", gin.H{
		"email": "u@test.com", "username": "user1", "password": "password123",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	var reg map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &reg))
	assert.NotEmpty(t, reg["token"])

	// Login
	w = doJSON(r, http.MethodPost, "/login", gin.H{
		"email": "u@test.com", "password": "password123",
	})
	require.Equal(t, http.StatusOK, w.Code)

	// Wrong password
	w = doJSON(r, http.MethodPost, "/login", gin.H{
		"email": "u@test.com", "password": "bad",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RegisterValidation(t *testing.T) {
	r := setupAuthRouter()
	w := doJSON(r, http.MethodPost, "/register", gin.H{
		"email": "not-an-email", "username": "x", "password": "123",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
