package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/streampulse/backend/internal/domain/entity"
	"github.com/streampulse/backend/internal/domain/repository"
)

// FileStorage abstracts where uploaded audio files are persisted.
// The local implementation lives in infrastructure/storage; an S3
// implementation can be added without changing the usecase.
type FileStorage interface {
	Save(ctx context.Context, key string, src io.Reader) (publicURL string, err error)
	Delete(ctx context.Context, key string) error
}

var (
	ErrUnsupportedMedia = errors.New("unsupported media type")
	ErrFileTooLarge     = errors.New("file too large")
)

// AllowedAudioMimes lists the audio MIME types accepted on upload.
var AllowedAudioMimes = map[string]string{
	"audio/mpeg":  ".mp3",
	"audio/mp3":   ".mp3",
	"audio/aac":   ".aac",
	"audio/x-aac": ".aac",
	"audio/ogg":   ".ogg",
	"audio/wav":   ".wav",
	"audio/x-wav": ".wav",
}

const MaxUploadBytes = 50 * 1024 * 1024 // 50 MB

type TrackUseCase struct {
	repo    repository.TrackRepository
	storage FileStorage
}

func NewTrackUseCase(repo repository.TrackRepository, storage FileStorage) *TrackUseCase {
	return &TrackUseCase{repo: repo, storage: storage}
}

// Upload stores an audio file, persists the metadata and returns the entity.
// The caller must have already validated that uploaderID has the broadcaster
// role at the HTTP layer.
func (uc *TrackUseCase) Upload(
	ctx context.Context,
	uploaderID, title, artist, contentType string,
	src io.Reader,
) (*entity.Track, error) {
	ext, ok := AllowedAudioMimes[contentType]
	if !ok {
		return nil, ErrUnsupportedMedia
	}
	key := uuid.NewString() + ext

	url, err := uc.storage.Save(ctx, key, src)
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	t := &entity.Track{
		Title:    title,
		Artist:   artist,
		FileURL:  url,
		UploadBy: uploaderID,
	}
	if err := uc.repo.Create(ctx, t); err != nil {
		// best-effort cleanup of orphaned file
		_ = uc.storage.Delete(context.Background(), key)
		return nil, fmt.Errorf("persist track: %w", err)
	}
	return t, nil
}

func (uc *TrackUseCase) List(ctx context.Context, offset, limit int) ([]entity.Track, error) {
	return uc.repo.List(ctx, offset, limit)
}

func (uc *TrackUseCase) Get(ctx context.Context, id string) (*entity.Track, error) {
	return uc.repo.FindByID(ctx, id)
}

// Delete removes a track. Only the original uploader (or an admin) may delete.
func (uc *TrackUseCase) Delete(ctx context.Context, id, callerID, callerRole string) error {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if t.UploadBy != callerID && callerRole != string(entity.RoleAdmin) {
		return ErrForbidden
	}
	return uc.repo.Delete(ctx, id)
}

// LocalStorage stores files on the local filesystem. It is the default for
// development; production should use an S3-compatible backend.
type LocalStorage struct {
	BasePath  string // e.g. ./uploads
	PublicURL string // e.g. http://localhost:8080/uploads (prefix used to build the URL)
}

func NewLocalStorage(basePath, publicURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir storage: %w", err)
	}
	return &LocalStorage{BasePath: basePath, PublicURL: publicURL}, nil
}

func (s *LocalStorage) Save(_ context.Context, key string, src io.Reader) (string, error) {
	target := filepath.Join(s.BasePath, key)
	out, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, io.LimitReader(src, MaxUploadBytes+1)); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return s.PublicURL + "/" + key, nil
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(s.BasePath, key))
}
