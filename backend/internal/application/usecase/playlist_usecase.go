package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

type PlaylistUseCase struct {
	playlists repository.PlaylistRepository
	tracks    repository.TrackRepository
}

func NewPlaylistUseCase(playlists repository.PlaylistRepository, tracks repository.TrackRepository) *PlaylistUseCase {
	return &PlaylistUseCase{playlists: playlists, tracks: tracks}
}

func (uc *PlaylistUseCase) Create(ctx context.Context, ownerID, name, description string) (*entity.Playlist, error) {
	p := &entity.Playlist{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	if err := uc.playlists.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create playlist: %w", err)
	}
	return p, nil
}

func (uc *PlaylistUseCase) Get(ctx context.Context, id string) (*entity.Playlist, error) {
	return uc.playlists.FindByID(ctx, id)
}

func (uc *PlaylistUseCase) ListByOwner(ctx context.Context, ownerID string) ([]entity.Playlist, error) {
	return uc.playlists.ListByOwner(ctx, ownerID)
}

func (uc *PlaylistUseCase) Update(ctx context.Context, id, callerID, name, description string) (*entity.Playlist, error) {
	p, err := uc.playlists.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.OwnerID != callerID {
		return nil, ErrForbidden
	}
	if name != "" {
		p.Name = name
	}
	p.Description = description
	if err := uc.playlists.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return p, nil
}

func (uc *PlaylistUseCase) Delete(ctx context.Context, id, callerID string) error {
	p, err := uc.playlists.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p.OwnerID != callerID {
		return ErrForbidden
	}
	return uc.playlists.Delete(ctx, id)
}

func (uc *PlaylistUseCase) AddTrack(ctx context.Context, playlistID, trackID, callerID string) error {
	p, err := uc.playlists.FindByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.OwnerID != callerID {
		return ErrForbidden
	}
	if _, err := uc.tracks.FindByID(ctx, trackID); err != nil {
		return err
	}
	return uc.playlists.AddTrack(ctx, playlistID, trackID)
}

func (uc *PlaylistUseCase) RemoveTrack(ctx context.Context, playlistID, trackID, callerID string) error {
	p, err := uc.playlists.FindByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.OwnerID != callerID {
		return ErrForbidden
	}
	return uc.playlists.RemoveTrack(ctx, playlistID, trackID)
}

func (uc *PlaylistUseCase) ReorderTracks(ctx context.Context, playlistID string, orderedTrackIDs []string, callerID string) error {
	p, err := uc.playlists.FindByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.OwnerID != callerID {
		return ErrForbidden
	}
	return uc.playlists.ReorderTracks(ctx, playlistID, orderedTrackIDs)
}
