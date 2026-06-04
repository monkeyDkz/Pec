package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/application/usecase"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/infrastructure/auth"
	"github.com/streampulse/backend/internal/infrastructure/config"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/streampulse/backend/internal/infrastructure/streaming"
	"github.com/streampulse/backend/internal/transport/http/handler"
	"github.com/streampulse/backend/internal/transport/http/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- in-memory repositories ----

type memUserRepo struct{ byID, byEmail map[string]*entity.User }

func newMemUserRepo() *memUserRepo {
	return &memUserRepo{byID: map[string]*entity.User{}, byEmail: map[string]*entity.User{}}
}
func (r *memUserRepo) Create(_ context.Context, u *entity.User) error {
	if _, ok := r.byEmail[u.Email]; ok {
		return persistence.ErrNotFound // simulate conflict path
	}
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	r.byID[u.ID], r.byEmail[u.Email] = u, u
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
func (r *memUserRepo) Update(_ context.Context, u *entity.User) error {
	r.byID[u.ID], r.byEmail[u.Email] = u, u
	return nil
}
func (r *memUserRepo) Delete(_ context.Context, id string) error { delete(r.byID, id); return nil }
func (r *memUserRepo) List(_ context.Context, _, _ int) ([]entity.User, error) {
	out := []entity.User{}
	for _, u := range r.byID {
		out = append(out, *u)
	}
	return out, nil
}

type memStreamRepo struct{ streams map[string]*entity.Stream }

func newMemStreamRepo() *memStreamRepo { return &memStreamRepo{streams: map[string]*entity.Stream{}} }
func (r *memStreamRepo) Create(_ context.Context, s *entity.Stream) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	r.streams[s.ID] = s
	return nil
}
func (r *memStreamRepo) FindByID(_ context.Context, id string) (*entity.Stream, error) {
	if s, ok := r.streams[id]; ok {
		return s, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *memStreamRepo) Update(_ context.Context, s *entity.Stream) error {
	r.streams[s.ID] = s
	return nil
}
func (r *memStreamRepo) Delete(_ context.Context, id string) error { delete(r.streams, id); return nil }
func (r *memStreamRepo) ListLive(_ context.Context) ([]entity.Stream, error) {
	out := []entity.Stream{}
	for _, s := range r.streams {
		if s.Status == entity.StreamStatusLive {
			out = append(out, *s)
		}
	}
	return out, nil
}
func (r *memStreamRepo) ListByBroadcaster(_ context.Context, _ string) ([]entity.Stream, error) {
	return nil, nil
}

type memPlaylistRepo struct{ pl map[string]*entity.Playlist }

func newMemPlaylistRepo() *memPlaylistRepo { return &memPlaylistRepo{pl: map[string]*entity.Playlist{}} }
func (r *memPlaylistRepo) Create(_ context.Context, p *entity.Playlist) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	r.pl[p.ID] = p
	return nil
}
func (r *memPlaylistRepo) FindByID(_ context.Context, id string) (*entity.Playlist, error) {
	if p, ok := r.pl[id]; ok {
		return p, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *memPlaylistRepo) Update(_ context.Context, p *entity.Playlist) error { r.pl[p.ID] = p; return nil }
func (r *memPlaylistRepo) Delete(_ context.Context, id string) error          { delete(r.pl, id); return nil }
func (r *memPlaylistRepo) ListByOwner(_ context.Context, owner string) ([]entity.Playlist, error) {
	out := []entity.Playlist{}
	for _, p := range r.pl {
		if p.OwnerID == owner {
			out = append(out, *p)
		}
	}
	return out, nil
}
func (r *memPlaylistRepo) AddTrack(_ context.Context, pid, tid string) error {
	r.pl[pid].Tracks = append(r.pl[pid].Tracks, entity.Track{ID: tid})
	return nil
}
func (r *memPlaylistRepo) RemoveTrack(_ context.Context, pid, tid string) error {
	p := r.pl[pid]
	kept := p.Tracks[:0]
	for _, t := range p.Tracks {
		if t.ID != tid {
			kept = append(kept, t)
		}
	}
	p.Tracks = kept
	return nil
}

type memTrackRepo struct{ tr map[string]*entity.Track }

func newMemTrackRepo() *memTrackRepo { return &memTrackRepo{tr: map[string]*entity.Track{}} }
func (r *memTrackRepo) Create(_ context.Context, t *entity.Track) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	r.tr[t.ID] = t
	return nil
}
func (r *memTrackRepo) FindByID(_ context.Context, id string) (*entity.Track, error) {
	if t, ok := r.tr[id]; ok {
		return t, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *memTrackRepo) Delete(_ context.Context, id string) error { delete(r.tr, id); return nil }
func (r *memTrackRepo) List(_ context.Context, _, _ int) ([]entity.Track, error) {
	out := []entity.Track{}
	for _, t := range r.tr {
		out = append(out, *t)
	}
	return out, nil
}

type memStatsRepo struct{ users *memUserRepo }

func (r memStatsRepo) Gather(_ context.Context) (repository.Stats, error) {
	return repository.Stats{TotalUsers: int64(len(r.users.byID))}, nil
}

// ---- test harness ----

type harness struct {
	r          http.Handler
	jwt        *auth.JWTManager
	users      *memUserRepo
}

func newHarness() *harness {
	cfg := &config.Config{
		JWTSecret:   "test-secret",
		JWTDuration: time.Hour,
		Environment: "test",
		CORSOrigins: "*",
	}
	users := newMemUserRepo()
	streams := newMemStreamRepo()
	playlists := newMemPlaylistRepo()
	tracks := newMemTrackRepo()
	stats := memStatsRepo{users: users}

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)
	hasher := auth.NewBcryptHasher()
	engine := streaming.NewEngine(streams)

	handlers := &router.Handlers{
		Auth:     handler.NewAuthHandler(usecase.NewAuthUseCase(users, jwtMgr, hasher)),
		User:     handler.NewUserHandler(usecase.NewUserUseCase(users)),
		Stream:   handler.NewStreamHandler(usecase.NewStreamUseCase(engine, streams)),
		Playlist: handler.NewPlaylistHandler(usecase.NewPlaylistUseCase(playlists, tracks)),
		Track:    handler.NewTrackHandler(usecase.NewTrackUseCase(tracks)),
		Admin:    handler.NewAdminHandler(usecase.NewAdminUseCase(users, stats)),
	}
	return &harness{r: router.New(cfg, handlers), jwt: jwtMgr, users: users}
}

func (h *harness) do(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.r.ServeHTTP(w, req)
	return w
}

// adminToken seeds an admin user and returns a valid JWT for it.
func (h *harness) adminToken(t *testing.T) string {
	return h.tokenForRole(t, string(entity.RoleAdmin))
}

// tokenForRole seeds a user with the given role and returns a valid JWT.
func (h *harness) tokenForRole(t *testing.T, role string) string {
	t.Helper()
	u := &entity.User{
		Email:    uuid.NewString() + "@x.com",
		Username: "u-" + uuid.NewString()[:8],
		Role:     entity.Role(role),
	}
	require.NoError(t, h.users.Create(context.Background(), u))
	tok, err := h.jwt.Generate(u.ID, role)
	require.NoError(t, err)
	return tok
}

func decode(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &m))
	return m
}

// ---- tests ----

func TestE2E_AuthAndProfile(t *testing.T) {
	h := newHarness()

	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	token, _ := decode(t, w)["token"].(string)
	require.NotEmpty(t, token)

	// unauthenticated /me is rejected
	w = h.do(t, http.MethodGet, "/api/v1/users/me", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// authenticated /me works
	w = h.do(t, http.MethodGet, "/api/v1/users/me", token, nil)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "user1", decode(t, w)["username"])

	// update profile
	w = h.do(t, http.MethodPut, "/api/v1/users/me", token, map[string]string{"username": "renamed"})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "renamed", decode(t, w)["username"])
}

func TestE2E_StreamRoleGating(t *testing.T) {
	h := newHarness()

	// a normal user
	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	userTok := decode(t, w)["token"].(string)

	// user cannot create a stream
	w = h.do(t, http.MethodPost, "/api/v1/streams", userTok, map[string]string{"title": "x"})
	assert.Equal(t, http.StatusForbidden, w.Code)

	// admin can
	adminTok := h.adminToken(t)
	w = h.do(t, http.MethodPost, "/api/v1/streams", adminTok, map[string]string{"title": "Radio", "description": "live"})
	require.Equal(t, http.StatusCreated, w.Code)
	sid := decode(t, w)["id"].(string)
	require.NotEmpty(t, sid)

	// it appears in the live list
	w = h.do(t, http.MethodGet, "/api/v1/streams", userTok, nil)
	require.Equal(t, http.StatusOK, w.Code)
	var list []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Len(t, list, 1)

	// get one
	w = h.do(t, http.MethodGet, "/api/v1/streams/"+sid, userTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// stop + delete by admin
	w = h.do(t, http.MethodPost, "/api/v1/streams/"+sid+"/stop", adminTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	w = h.do(t, http.MethodDelete, "/api/v1/streams/"+sid, adminTok, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestE2E_PlaylistsAndTracks(t *testing.T) {
	h := newHarness()
	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	userTok := decode(t, w)["token"].(string)
	adminTok := h.adminToken(t)

	// admin/broadcaster creates a track
	w = h.do(t, http.MethodPost, "/api/v1/tracks", adminTok, map[string]any{
		"title": "Song", "artist": "Artist", "duration": 180, "file_url": "https://cdn/x.mp3",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	tid := decode(t, w)["id"].(string)

	// a normal user cannot create tracks
	w = h.do(t, http.MethodPost, "/api/v1/tracks", userTok, map[string]any{
		"title": "x", "file_url": "https://cdn/y.mp3",
	})
	assert.Equal(t, http.StatusForbidden, w.Code)

	// user creates a playlist and adds the track
	w = h.do(t, http.MethodPost, "/api/v1/playlists", userTok, map[string]string{"name": "Chill"})
	require.Equal(t, http.StatusCreated, w.Code)
	pid := decode(t, w)["id"].(string)

	w = h.do(t, http.MethodPost, "/api/v1/playlists/"+pid+"/tracks", userTok, map[string]string{"track_id": tid})
	assert.Equal(t, http.StatusOK, w.Code)

	// playlist now has the track
	w = h.do(t, http.MethodGet, "/api/v1/playlists/"+pid, userTok, nil)
	require.Equal(t, http.StatusOK, w.Code)
	tracks := decode(t, w)["tracks"].([]any)
	assert.Len(t, tracks, 1)

	// list user's playlists
	w = h.do(t, http.MethodGet, "/api/v1/playlists", userTok, nil)
	require.Equal(t, http.StatusOK, w.Code)

	// remove track + delete playlist
	w = h.do(t, http.MethodDelete, "/api/v1/playlists/"+pid+"/tracks/"+tid, userTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	w = h.do(t, http.MethodDelete, "/api/v1/playlists/"+pid, userTok, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestE2E_AdminEndpoints(t *testing.T) {
	h := newHarness()
	// register a normal user, then build an admin
	_ = h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	adminTok := h.adminToken(t)

	// non-admin token blocked
	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "v@x.com", "username": "user2", "password": "password123",
	})
	userTok := decode(t, w)["token"].(string)
	w = h.do(t, http.MethodGet, "/api/v1/admin/stats", userTok, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// admin stats
	w = h.do(t, http.MethodGet, "/api/v1/admin/stats", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, decode(t, w)["total_users"], float64(1))

	// admin lists users
	w = h.do(t, http.MethodGet, "/api/v1/admin/users", adminTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}
