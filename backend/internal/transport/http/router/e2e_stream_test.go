package router_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_ListenStreamsLiveAudio drives the real streaming handler over an HTTP
// server: a listener connects, the broadcaster publishes, and the bytes arrive.
func TestE2E_ListenStreamsLiveAudio(t *testing.T) {
	h := newHarness()
	srv := httptest.NewServer(h.r)
	defer srv.Close()
	adminTok := h.adminToken(t)

	// Create a live stream over real HTTP.
	body := strings.NewReader(`{"title":"Live"}`)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/streams", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminTok)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	var created map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	resp.Body.Close()
	sid := created["id"].(string)

	// Connect a listener; the handler subscribes synchronously on entry.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lreq, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/v1/streams/"+sid+"/listen", nil)
	lreq.Header.Set("Authorization", "Bearer "+adminTok)

	// Publish after a short delay so the listener is subscribed first.
	go func() {
		time.Sleep(150 * time.Millisecond)
		preq, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/streams/"+sid+"/publish", strings.NewReader("LIVE_AUDIO_XYZ"))
		preq.Header.Set("Authorization", "Bearer "+adminTok)
		if pr, err := http.DefaultClient.Do(preq); err == nil {
			pr.Body.Close()
		}
	}()

	lresp, err := http.DefaultClient.Do(lreq)
	require.NoError(t, err)
	defer lresp.Body.Close()
	require.Equal(t, http.StatusOK, lresp.StatusCode)

	buf := make([]byte, len("LIVE_AUDIO_XYZ"))
	n, err := io.ReadFull(lresp.Body, buf)
	require.NoError(t, err)
	assert.Equal(t, "LIVE_AUDIO_XYZ", string(buf[:n]))
}
