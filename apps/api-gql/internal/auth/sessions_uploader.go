package auth

import (
	"context"
	"fmt"
)

const latestUploadedFilesIdsKey = "latestUploadedFilesIds"

func (s *Auth) GetLatestUploadedFilesIds(ctx context.Context) ([]string, error) {
	ids, ok := s.sessionManager.Get(ctx, latestUploadedFilesIdsKey).([]string)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	return ids, nil
}

func (s *Auth) AddLatestUploadedFileId(ctx context.Context, id string) error {
	latest, err := s.GetLatestUploadedFilesIds(ctx)
	if err != nil {
		latest = nil
	}
	latest = append([]string{id}, latest...)
	if len(latest) > 50 {
		latest = latest[:50]
	}

	s.sessionManager.Put(ctx, latestUploadedFilesIdsKey, latest)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}
