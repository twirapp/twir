package notifications_sync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	cfg "github.com/twirapp/twir/libs/config"
)

type mediaSource struct {
	ID          string
	URL         string
	Filename    string
	ContentType string
	Size        int64
}

type attachmentStore struct {
	client    *minio.Client
	http      *http.Client
	bucket    string
	publicURL string
	maxBytes  int64
}

var unsafeFilenamePattern = regexp.MustCompile("[^a-zA-Z0-9._-]+")

func newAttachmentStore(config cfg.Config) (*attachmentStore, error) {
	if config.S3Host == "" || config.S3Bucket == "" || config.S3PublicUrl == "" {
		return nil, nil
	}

	accessToken := config.S3AccessToken
	secretToken := config.S3SecretToken
	if config.AppEnv != "production" {
		accessToken = "minio"
		secretToken = "minio-password"
	}

	client, err := minio.New(config.S3Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessToken, secretToken, ""),
		Region: config.S3Region,
		Secure: config.AppEnv == "production",
	})
	if err != nil {
		return nil, fmt.Errorf("create notifications object store: %w", err)
	}

	maxBytes := config.DiscordNotificationsMaxAttachmentBytes
	if maxBytes <= 0 {
		maxBytes = 25 << 20
	}

	return &attachmentStore{
		client: client,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		bucket:    config.S3Bucket,
		publicURL: strings.TrimRight(config.S3PublicUrl, "/"),
		maxBytes:  maxBytes,
	}, nil
}

func mediaSources(message discord.Message) []mediaSource {
	result := make([]mediaSource, 0, len(message.Attachments)+len(message.Embeds)*2)
	seenURLs := make(map[string]struct{})
	appendSource := func(source mediaSource) {
		if source.URL == "" {
			return
		}
		if _, ok := seenURLs[source.URL]; ok {
			return
		}
		seenURLs[source.URL] = struct{}{}
		result = append(result, source)
	}

	for _, attachment := range message.Attachments {
		appendSource(mediaSource{
			ID:          attachment.ID.String(),
			URL:         string(attachment.URL),
			Filename:    attachment.Filename,
			ContentType: attachment.ContentType,
			Size:        int64(attachment.Size),
		})
	}

	for _, embed := range message.Embeds {
		if embed.Image != nil {
			appendSource(sourceFromEmbedURL(string(embed.Image.URL)))
		}
		if embed.Thumbnail != nil {
			appendSource(sourceFromEmbedURL(string(embed.Thumbnail.URL)))
		}
	}

	return result
}

func sourceFromEmbedURL(value string) mediaSource {
	sum := sha256.Sum256([]byte(value))
	filename := "embed-image"
	if parsed, err := url.Parse(value); err == nil {
		if base := path.Base(parsed.Path); base != "" && base != "." && base != "/" {
			filename = base
		}
	}

	return mediaSource{
		ID:       fmt.Sprintf("%x", sum[:8]),
		URL:      value,
		Filename: filename,
	}
}

func isImage(contentType string, filename string) bool {
	if strings.HasPrefix(strings.ToLower(contentType), "image/") {
		return true
	}
	guessed := mime.TypeByExtension(strings.ToLower(path.Ext(filename)))
	return strings.HasPrefix(guessed, "image/")
}

func safeFilename(value string) string {
	value = unsafeFilenamePattern.ReplaceAllString(path.Base(value), "_")
	value = strings.Trim(value, "._")
	if value == "" {
		return "image"
	}
	return value
}

func (s *attachmentStore) persist(
	ctx context.Context,
	messageID string,
	source mediaSource,
) (renderedMedia, string, error) {
	fallback := renderedMedia{
		URL:         source.URL,
		Filename:    source.Filename,
		ContentType: source.ContentType,
		IsImage:     isImage(source.ContentType, source.Filename),
	}
	if s == nil || !fallback.IsImage {
		return fallback, "", nil
	}
	if source.Size > s.maxBytes {
		return fallback, "", fmt.Errorf(
			"attachment %q is %d bytes, limit is %d",
			source.Filename,
			source.Size,
			s.maxBytes,
		)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
	if err != nil {
		return fallback, "", err
	}
	response, err := s.http.Do(request)
	if err != nil {
		return fallback, "", err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fallback, "", fmt.Errorf("download attachment: HTTP %d", response.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, s.maxBytes+1))
	if err != nil {
		return fallback, "", err
	}
	if int64(len(data)) > s.maxBytes {
		return fallback, "", fmt.Errorf("downloaded attachment exceeds %d bytes", s.maxBytes)
	}

	contentType := source.ContentType
	if contentType == "" {
		contentType = response.Header.Get("Content-Type")
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		return fallback, "", fmt.Errorf("attachment content type %q is not an image", contentType)
	}

	objectKey := fmt.Sprintf(
		"notifications/discord/%s/%s-%s",
		messageID,
		safeFilename(source.ID),
		safeFilename(source.Filename),
	)
	_, err = s.client.PutObject(
		ctx,
		s.bucket,
		objectKey,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fallback, "", fmt.Errorf("upload attachment: %w", err)
	}

	return renderedMedia{
		URL:         s.publicURL + "/" + objectKey,
		Filename:    source.Filename,
		ContentType: contentType,
		IsImage:     true,
	}, objectKey, nil
}

func (s *attachmentStore) remove(ctx context.Context, objectKeys []string) error {
	if s == nil {
		return nil
	}
	for _, objectKey := range objectKeys {
		if err := s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("remove attachment %q: %w", objectKey, err)
		}
	}
	return nil
}
