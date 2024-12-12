package resolvers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func (r *Resolver) computeBadgeFileName(file graphql.Upload, fileID uuid.UUID) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	if fileExtension == "" {
		return "", fmt.Errorf("file extension is empty")
	}
	if !strings.HasPrefix(file.ContentType, "image/") {
		return "", fmt.Errorf("file is not an image")
	}

	fileExtension = strings.ToLower(fileExtension)
	fileName := fmt.Sprintf("%s%s", fileID, fileExtension)

	return fileName, nil
}
