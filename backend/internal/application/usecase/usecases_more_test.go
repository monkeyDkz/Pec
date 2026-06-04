package usecase

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/domain/service"
	"github.com/streampulse/backend/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- fakes ----

type fakeEngine struct {
	listeners  map[string]int
	publishErr error
}

func (f *fakeEngine) StartStream(_ context.Context, b, t, d string) (*entity.Stream, error) {
	return &entity.Stream{ID: uuid.NewString(), Title: t, Description: d, BroadcasterID: b, Status: entity.StreamStatusLive}, nil
}
func (f *fakeEngine) StopStream(_ context.Context, _, _ string) error { return nil }
func (f *fakeEngine) ListenStream(_ context.Context, _ string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}
func (f *fakeEngine) PublishAudio(_ context.Context, _ string, _ io.Reader) error {
	return f.publishErr
}
func (f *fakeEngine) GetListenerCount(_ context.Context, id string) (int, error) {
	if f.listeners == nil {
		return 0, nil
	}
	return f.listeners[id], nil
}

type fakeStreamRepoUC struct {
	streams map[string]*entity.Stream
}

func newFakeStreamRepoUC() *fakeStreamRepoUC {
	return &fakeStreamRepoUC{streams: map[string]*entity.Stream{}}
}
func (r *fakeStreamRepoUC) Create(_ context.Context, s *entity.Stream) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	r.streams[s.ID] = s
	return nil
}
func (r *fakeStreamRepoUC) FindByID(_ context.Context, id string) (*entity.Stream, error) {
	if s, ok := r.streams[id]; ok {
		return s, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *fakeStreamRepoUC) Update(_ context.Context, s *entity.Stream) error {
	r.streams[s.ID] = s
	return nil
}
func (r *fakeStreamRepoUC) Delete(_ context.Context, id string) error {
	delete(r.streams, id)
	return nil
}
func (r *fakeStreamRepoUC) ListLive(_ context.Context) ([]entity.Stream, error) {
	out := []entity.Stream{}
	for _, s := range r.streams {
		if s.Status == entity.StreamStatusLive {
			out = append(out, *s)
		}
	}
	return out, nil
}
func (r *fakeStreamRepoUC) ListByBroadcaster(_ context.Context, b string) ([]entity.Stream, error) {
	return nil, nil
}

type fakeTrackRepo struct{ tracks map[string]*entity.Track }

func newFakeTrackRepo() *fakeTrackRepo { return &fakeTrackRepo{tracks: map[string]*entity.Track{}} }
func (r *fakeTrackRepo) Create(_ context.Context, t *entity.Track) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	r.tracks[t.ID] = t
	return nil
}
func (r *fakeTrackRepo) FindByID(_ context.Context, id string) (*entity.Track, error) {
	if t, ok := r.tracks[id]; ok {
		return t, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *fakeTrackRepo) Delete(_ context.Context, id string) error { delete(r.tracks, id); return nil }
func (r *fakeTrackRepo) List(_ context.Context, _, _ int) ([]entity.Track, error) {
	out := []entity.Track{}
	for _, t := range r.tracks {
		out = append(out, *t)
	}
	return out, nil
}

type fakePlaylistRepo struct{ playlists map[string]*entity.Playlist }

func newFakePlaylistRepo() *fakePlaylistRepo {
	return &fakePlaylistRepo{playlists: map[string]*entity.Playlist{}}
}
func (r *fakePlaylistRepo) Create(_ context.Context, p *entity.Playlist) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	r.playlists[p.ID] = p
	return nil
}
func (r *fakePlaylistRepo) FindByID(_ context.Context, id string) (*entity.Playlist, error) {
	if p, ok := r.playlists[id]; ok {
		return p, nil
	}
	return nil, persistence.ErrNotFound
}
func (r *fakePlaylistRepo) Update(_ context.Context, p *entity.Playlist) error {
	r.playlists[p.ID] = p
	return nil
}
func (r *fakePlaylistRepo) Delete(_ context.Context, id string) error {
	delete(r.playlists, id)
	return nil
}
func (r *fakePlaylistRepo) ListByOwner(_ context.Context, owner string) ([]entity.Playlist, error) {
	out := []entity.Playlist{}
	for _, p := range r.playlists {
		if p.OwnerID == owner {
			out = append(out, *p)
		}
	}
	return out, nil
}
func (r *fakePlaylistRepo) AddTrack(_ context.Context, pid, tid string) error {
	p := r.playlists[pid]
	p.Tracks = append(p.Tracks, entity.Track{ID: tid})
	return nil
}
func (r *fakePlaylistRepo) RemoveTrack(_ context.Context, pid, tid string) error {
	p := r.playlists[pid]
	kept := p.Tracks[:0]
	for _, t := range p.Tracks {
		if t.ID != tid {
			kept = append(kept, t)
		}
	}
	p.Tracks = kept
	return nil
}

type fakeStatsRepo struct{}

func (fakeStatsRepo) Gather(_ context.Context) (repository.Stats, error) {
	return repository.Stats{TotalUsers: 5, LiveStreams: 2}, nil
}

// ---- tests ----

func TestStreamUseCase_PublishOwnership(t *testing.T) {
	ctx := context.Background()
	repo := newFakeStreamRepoUC()
	s := &entity.Stream{BroadcasterID: "owner", Status: entity.StreamStatusLive}
	_ = repo.Create(ctx, s)
	uc := NewStreamUseCase(&fakeEngine{}, repo)

	err := uc.Publish(ctx, s.ID, "intruder", strings.NewReader("x"))
	assert.ErrorIs(t, err, service.ErrForbidden)

	err = uc.Publish(ctx, s.ID, "owner", strings.NewReader("x"))
	assert.NoError(t, err)
}

func TestStreamUseCase_ListLiveFillsListenerCount(t *testing.T) {
	ctx := context.Background()
	repo := newFakeStreamRepoUC()
	s := &entity.Stream{BroadcasterID: "o", Status: entity.StreamStatusLive}
	_ = repo.Create(ctx, s)
	uc := NewStreamUseCase(&fakeEngine{listeners: map[string]int{s.ID: 7}}, repo)

	live, err := uc.ListLive(ctx)
	require.NoError(t, err)
	require.Len(t, live, 1)
	assert.Equal(t, 7, live[0].ListenerCount)
}

func TestStreamUseCase_DeletePermissions(t *testing.T) {
	ctx := context.Background()
	repo := newFakeStreamRepoUC()
	s := &entity.Stream{BroadcasterID: "owner", Status: entity.StreamStatusLive}
	_ = repo.Create(ctx, s)
	uc := NewStreamUseCase(&fakeEngine{}, repo)

	assert.ErrorIs(t, uc.Delete(ctx, s.ID, "other", "user"), service.ErrForbidden)
	assert.NoError(t, uc.Delete(ctx, s.ID, "other", "admin")) // admin override
}

func TestPlaylistUseCase_Lifecycle(t *testing.T) {
	ctx := context.Background()
	prepo := newFakePlaylistRepo()
	trepo := newFakeTrackRepo()
	uc := NewPlaylistUseCase(prepo, trepo)

	pl, err := uc.Create(ctx, "owner", "Chill", "desc")
	require.NoError(t, err)
	assert.Equal(t, "Chill", pl.Name)

	track := &entity.Track{Title: "Song"}
	_ = trepo.Create(ctx, track)

	// non-owner cannot add
	assert.ErrorIs(t, uc.AddTrack(ctx, pl.ID, "intruder", track.ID), service.ErrForbidden)
	// owner can add
	require.NoError(t, uc.AddTrack(ctx, pl.ID, "owner", track.ID))

	got, _ := uc.Get(ctx, pl.ID)
	require.Len(t, got.Tracks, 1)

	require.NoError(t, uc.RemoveTrack(ctx, pl.ID, "owner", track.ID))
	got, _ = uc.Get(ctx, pl.ID)
	assert.Len(t, got.Tracks, 0)

	// update + delete guarded by ownership
	_, err = uc.Update(ctx, pl.ID, "intruder", "New", "")
	assert.ErrorIs(t, err, service.ErrForbidden)
	assert.ErrorIs(t, uc.Delete(ctx, pl.ID, "intruder"), service.ErrForbidden)
	require.NoError(t, uc.Delete(ctx, pl.ID, "owner"))
}

func TestTrackUseCase_CreateAndList(t *testing.T) {
	ctx := context.Background()
	uc := NewTrackUseCase(newFakeTrackRepo())

	tr, err := uc.Create(ctx, "uploader", "Title", "Artist", 180, "https://cdn/x.mp3")
	require.NoError(t, err)
	assert.Equal(t, "uploader", tr.UploadBy)

	list, err := uc.List(ctx, 0, 0) // 0 limit -> defaulted
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestUserUseCase_UpdateUsername(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()
	u := &entity.User{Username: "old", Email: "e@x.com"}
	_ = repo.Create(ctx, u)
	uc := NewUserUseCase(repo)

	got, err := uc.UpdateUsername(ctx, u.ID, "new")
	require.NoError(t, err)
	assert.Equal(t, "new", got.Username)
}

func TestAdminUseCase_RoleAndStats(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()
	u := &entity.User{Username: "u", Email: "e@x.com", Role: entity.RoleUser}
	_ = repo.Create(ctx, u)
	uc := NewAdminUseCase(repo, fakeStatsRepo{})

	got, err := uc.UpdateRole(ctx, u.ID, "broadcaster")
	require.NoError(t, err)
	assert.Equal(t, entity.RoleBroadcaster, got.Role)

	stats, err := uc.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(5), stats.TotalUsers)
}
