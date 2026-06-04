package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
	"github.com/streampulse/backend/internal/domain/service"
)

// PlaylistUseCase manages playlists and their track membership, enforcing ownership.
type PlaylistUseCase struct {
	repo      repository.PlaylistRepository
	trackRepo repository.TrackRepository
}

func NewPlaylistUseCase(repo repository.PlaylistRepository, trackRepo repository.TrackRepository) *PlaylistUseCase {
	return &PlaylistUseCase{repo: repo, trackRepo: trackRepo}
}

func (uc *PlaylistUseCase) Create(ctx context.Context, ownerID, name, description string) (*entity.Playlist, error) {
	playlist := &entity.Playlist{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	if err := uc.repo.Create(ctx, playlist); err != nil {
		return nil, fmt.Errorf("create playlist: %w", err)
	}
	return playlist, nil
}

func (uc *PlaylistUseCase) Get(ctx context.Context, id string) (*entity.Playlist, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *PlaylistUseCase) ListByOwner(ctx context.Context, ownerID string) ([]entity.Playlist, error) {
	return uc.repo.ListByOwner(ctx, ownerID)
}

func (uc *PlaylistUseCase) Update(ctx context.Context, id, ownerID, name, description string) (*entity.Playlist, error) {
	playlist, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find playlist: %w", err)
	}
	if playlist.OwnerID != ownerID {
		return nil, service.ErrForbidden
	}
	if name != "" {
		playlist.Name = name
	}
	if description != "" {
		playlist.Description = description
	}
	if err := uc.repo.Update(ctx, playlist); err != nil {
		return nil, fmt.Errorf("update playlist: %w", err)
	}
	return playlist, nil
}

func (uc *PlaylistUseCase) Delete(ctx context.Context, id, ownerID string) error {
	playlist, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find playlist: %w", err)
	}
	if playlist.OwnerID != ownerID {
		return service.ErrForbidden
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *PlaylistUseCase) AddTrack(ctx context.Context, playlistID, ownerID, trackID string) error {
	playlist, err := uc.repo.FindByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("find playlist: %w", err)
	}
	if playlist.OwnerID != ownerID {
		return service.ErrForbidden
	}
	if _, err := uc.trackRepo.FindByID(ctx, trackID); err != nil {
		return fmt.Errorf("find track: %w", err)
	}
	return uc.repo.AddTrack(ctx, playlistID, trackID)
}

func (uc *PlaylistUseCase) RemoveTrack(ctx context.Context, playlistID, ownerID, trackID string) error {
	playlist, err := uc.repo.FindByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("find playlist: %w", err)
	}
	if playlist.OwnerID != ownerID {
		return service.ErrForbidden
	}
	return uc.repo.RemoveTrack(ctx, playlistID, trackID)
}
