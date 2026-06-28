package channels_secret

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/repositories/channels_secret"
	"github.com/twirapp/twir/libs/repositories/channels_secret/model"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config            config.Config
	AuditRecorder     audit.Recorder
	SecretsRepository channels_secret.Repository
}

type Service struct {
	config            config.Config
	auditRecorder     audit.Recorder
	secretsRepository channels_secret.Repository
	encryptionKey     []byte
}

var ErrNotFound = errors.New("secret not found")

func New(opts Opts) (*Service, error) {
	key := opts.Config.SecretsEncryptionKey
	if key == "" {
		return nil, fmt.Errorf("SECRETS_ENCRYPTION_KEY is required")
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("SECRETS_ENCRYPTION_KEY must be 32 bytes, got %d", len(key))
	}

	return &Service{
		config:            opts.Config,
		auditRecorder:     opts.AuditRecorder,
		secretsRepository: opts.SecretsRepository,
		encryptionKey:     []byte(key),
	}, nil
}

func (s *Service) GetAllByChannelID(ctx context.Context, channelID string) (
	[]model.ChannelSecret,
	error,
) {
	return s.secretsRepository.GetAllByChannelID(ctx, channelID)
}

func (s *Service) GetAllDecryptedByChannelID(ctx context.Context, channelID string) (
	[]model.ChannelSecret,
	error,
) {
	secrets, err := s.secretsRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	for i, secret := range secrets {
		decrypted, err := s.decrypt(secret.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot decrypt secret %s: %w", secret.Name, err)
		}
		secrets[i].Value = decrypted
	}

	return secrets, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (
	model.ChannelSecret,
	error,
) {
	return s.secretsRepository.GetByID(ctx, id)
}

func (s *Service) GetDecryptedValue(ctx context.Context, id uuid.UUID) (string, error) {
	secret, err := s.secretsRepository.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	if secret.IsNil() {
		return "", ErrNotFound
	}

	decrypted, err := s.decrypt(secret.Value)
	if err != nil {
		return "", fmt.Errorf("cannot decrypt secret: %w", err)
	}

	return decrypted, nil
}

func (s *Service) GetDecryptedByChannelID(ctx context.Context, channelID string) (
	map[string]string,
	error,
) {
	secrets, err := s.secretsRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(secrets))
	for _, secret := range secrets {
		decrypted, err := s.decrypt(secret.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot decrypt secret %s: %w", secret.Name, err)
		}
		result[secret.Name] = decrypted
	}

	return result, nil
}

type CreateInput struct {
	ChannelID   string
	ActorID     string
	Name        string
	Description *string
	Value       string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (
	model.ChannelSecret,
	error,
) {
	encrypted, err := s.encrypt(input.Value)
	if err != nil {
		return model.Nil, fmt.Errorf("cannot encrypt secret: %w", err)
	}

	secret, err := s.secretsRepository.Create(ctx, channels_secret.CreateInput{
		ChannelID:   input.ChannelID,
		Name:        input.Name,
		Description: input.Description,
		Value:       encrypted,
	})
	if err != nil {
		return model.Nil, err
	}

	return secret, nil
}

type UpdateInput struct {
	Name        *string
	Description *string
	Value       *string
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	model.ChannelSecret,
	error,
) {
	repoInput := channels_secret.UpdateInput{
		Name:        input.Name,
		Description: input.Description,
	}

	if input.Value != nil {
		encrypted, err := s.encrypt(*input.Value)
		if err != nil {
			return model.Nil, fmt.Errorf("cannot encrypt secret: %w", err)
		}
		repoInput.Value = &encrypted
	}

	secret, err := s.secretsRepository.Update(ctx, id, repoInput)
	if err != nil {
		return model.Nil, err
	}

	return secret, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.secretsRepository.Delete(ctx, id)
}

func (s *Service) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Service) decrypt(encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
