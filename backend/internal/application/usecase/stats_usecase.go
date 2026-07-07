package usecase

import (
	"context"
	"fmt"

	"github.com/streampulse/backend/internal/domain/entity"
	"gorm.io/gorm"
)

// Stats is a snapshot of the platform usage, returned to admins.
type Stats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalBroadcasters int64 `json:"total_broadcasters"`
	TotalStreams     int64 `json:"total_streams"`
	LiveStreams      int64 `json:"live_streams"`
	TotalPlaylists   int64 `json:"total_playlists"`
	TotalTracks      int64 `json:"total_tracks"`
	TotalFeedbacks   int64 `json:"total_feedbacks"`
}

type StatsUseCase struct {
	db *gorm.DB
}

func NewStatsUseCase(db *gorm.DB) *StatsUseCase {
	return &StatsUseCase{db: db}
}

// Snapshot runs a small set of count queries in parallel transactions.
// Each is bounded by ctx.
func (uc *StatsUseCase) Snapshot(ctx context.Context) (*Stats, error) {
	var s Stats
	queries := []struct {
		target *int64
		query  func() *gorm.DB
	}{
		{&s.TotalUsers, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.User{})
		}},
		{&s.TotalBroadcasters, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.User{}).
				Where("role IN ?", []string{string(entity.RoleBroadcaster), string(entity.RoleAdmin)})
		}},
		{&s.TotalStreams, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.Stream{})
		}},
		{&s.LiveStreams, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.Stream{}).
				Where("status = ?", entity.StreamStatusLive)
		}},
		{&s.TotalPlaylists, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.Playlist{})
		}},
		{&s.TotalTracks, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.Track{})
		}},
		{&s.TotalFeedbacks, func() *gorm.DB {
			return uc.db.WithContext(ctx).Model(&entity.Feedback{})
		}},
	}

	for _, q := range queries {
		if err := q.query().Count(q.target).Error; err != nil {
			return nil, fmt.Errorf("count: %w", err)
		}
	}
	return &s, nil
}
