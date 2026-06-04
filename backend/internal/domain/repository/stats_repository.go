package repository

import "context"

// Stats aggregates platform-wide counts for the admin dashboard.
type Stats struct {
	TotalUsers     int64
	TotalStreams   int64
	LiveStreams    int64
	TotalTracks    int64
	TotalPlaylists int64
}

type StatsRepository interface {
	Gather(ctx context.Context) (Stats, error)
}
