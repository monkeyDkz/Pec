package router_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_TrackCRUD(t *testing.T) {
	h := newHarness()
	adminTok := h.adminToken(t)

	w := h.do(t, http.MethodPost, "/api/v1/tracks", adminTok, map[string]any{
		"title": "Song", "artist": "A", "duration": 120, "file_url": "https://cdn/x.mp3",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	tid := decode(t, w)["id"].(string)

	w = h.do(t, http.MethodGet, "/api/v1/tracks/"+tid, adminTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = h.do(t, http.MethodGet, "/api/v1/tracks", adminTok, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = h.do(t, http.MethodDelete, "/api/v1/tracks/"+tid, adminTok, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestE2E_PublishAndListenErrors(t *testing.T) {
	h := newHarness()
	adminTok := h.adminToken(t)

	w := h.do(t, http.MethodPost, "/api/v1/streams", adminTok, map[string]string{"title": "Radio"})
	require.Equal(t, http.StatusCreated, w.Code)
	sid := decode(t, w)["id"].(string)

	// owner publishes a short payload -> 200
	w = h.do(t, http.MethodPost, "/api/v1/streams/"+sid+"/publish", adminTok, map[string]string{"audio": "chunk"})
	assert.Equal(t, http.StatusOK, w.Code)

	// listening to a non-live stream -> 409 (covers Listen handler + respondError)
	w = h.do(t, http.MethodGet, "/api/v1/streams/"+uuid.NewString()+"/listen", adminTok, nil)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestE2E_AdminUpdateRole(t *testing.T) {
	h := newHarness()
	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	uid := decode(t, w)["user"].(map[string]any)["id"].(string)
	adminTok := h.adminToken(t)

	w = h.do(t, http.MethodPut, "/api/v1/admin/users/"+uid+"/role", adminTok, map[string]string{"role": "broadcaster"})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "broadcaster", decode(t, w)["role"])
}

func TestE2E_NotFound(t *testing.T) {
	h := newHarness()
	tok := h.adminToken(t)
	missing := uuid.NewString()

	for _, path := range []string{
		"/api/v1/streams/" + missing,
		"/api/v1/playlists/" + missing,
		"/api/v1/tracks/" + missing,
	} {
		w := h.do(t, http.MethodGet, path, tok, nil)
		assert.Equal(t, http.StatusNotFound, w.Code, "GET %s should 404", path)
	}
}

func TestE2E_PublishForbiddenForNonOwner(t *testing.T) {
	h := newHarness()
	adminTok := h.adminToken(t)

	// admin creates a stream
	w := h.do(t, http.MethodPost, "/api/v1/streams", adminTok, map[string]string{"title": "Radio"})
	sid := decode(t, w)["id"].(string)

	// another broadcaster cannot publish to it
	other := h.tokenForRole(t, "broadcaster")
	w = h.do(t, http.MethodPost, "/api/v1/streams/"+sid+"/publish", other, map[string]string{"a": "b"})
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_PlaylistUpdate(t *testing.T) {
	h := newHarness()
	w := h.do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]string{
		"email": "u@x.com", "username": "user1", "password": "password123",
	})
	userTok := decode(t, w)["token"].(string)

	w = h.do(t, http.MethodPost, "/api/v1/playlists", userTok, map[string]string{"name": "Old"})
	pid := decode(t, w)["id"].(string)

	w = h.do(t, http.MethodPut, "/api/v1/playlists/"+pid, userTok, map[string]string{"name": "New", "description": "d"})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "New", decode(t, w)["name"])
}
