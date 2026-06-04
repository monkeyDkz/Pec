package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeUserRepo is an in-memory UserRepository for unit tests.
type fakeUserRepo struct {
	byID    map[string]*entity.User
	byEmail map[string]*entity.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byID:    make(map[string]*entity.User),
		byEmail: make(map[string]*entity.User),
	}
}

func (r *fakeUserRepo) Create(_ context.Context, u *entity.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.Role == "" {
		u.Role = entity.RoleUser
	}
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *fakeUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, persistence.ErrNotFound
}

func (r *fakeUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, persistence.ErrNotFound
}

func (r *fakeUserRepo) Update(_ context.Context, u *entity.User) error {
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *fakeUserRepo) Delete(_ context.Context, id string) error {
	if u, ok := r.byID[id]; ok {
		delete(r.byEmail, u.Email)
		delete(r.byID, id)
	}
	return nil
}

func (r *fakeUserRepo) List(_ context.Context, offset, limit int) ([]entity.User, error) {
	out := make([]entity.User, 0, len(r.byID))
	for _, u := range r.byID {
		out = append(out, *u)
	}
	return out, nil
}

func newAuthUC() *AuthUseCase {
	return NewAuthUseCase(
		newFakeUserRepo(),
		auth.NewJWTManager("test-secret", time.Hour),
		auth.NewBcryptHasher(),
	)
}

func TestAuthUseCase_RegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	uc := newAuthUC()

	user, token, err := uc.Register(ctx, "a@b.com", "alice", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, entity.RoleUser, user.Role)
	assert.NotEmpty(t, user.ID)

	// Login with correct credentials.
	tok, err := uc.Login(ctx, "a@b.com", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, tok)

	// ValidateToken resolves the user.
	got, err := uc.ValidateToken(ctx, tok)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestAuthUseCase_LoginFailures(t *testing.T) {
	ctx := context.Background()
	uc := newAuthUC()
	_, _, err := uc.Register(ctx, "a@b.com", "alice", "password123")
	require.NoError(t, err)

	t.Run("wrong password", func(t *testing.T) {
		_, err := uc.Login(ctx, "a@b.com", "nope")
		require.Error(t, err)
	})

	t.Run("unknown email", func(t *testing.T) {
		_, err := uc.Login(ctx, "ghost@b.com", "password123")
		require.Error(t, err)
	})
}
