package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

// TrackUseCase manages audio track sources.
type TrackUseCase struct {
	repo repository.TrackRepository
}

func NewTrackUseCase(repo repository.TrackRepository) *TrackUseCase {
	return &TrackUseCase{repo: repo}
}

func (uc *TrackUseCase) Create(ctx context.Context, uploaderID, title, artist string, duration int, fileURL string) (*entity.Track, error) {
	track := &entity.Track{
		Title:    title,
		Artist:   artist,
		Duration: duration,
		FileURL:  fileURL,
		UploadBy: uploaderID,
	}
	if err := uc.repo.Create(ctx, track); err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}
	return track, nil
}

func (uc *TrackUseCase) Get(ctx context.Context, id string) (*entity.Track, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *TrackUseCase) List(ctx context.Context, offset, limit int) ([]entity.Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return uc.repo.List(ctx, offset, limit)
}

func (uc *TrackUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
