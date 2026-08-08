package uploader

import (
	"errors"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	config "github.com/twirapp/twir/libs/config"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
)

type Registration struct{}

type registerRoute interface {
	Register(huma.API)
}

func RegisterRoutes(
	api huma.API,
	config config.Config,
	service *uploaderservice.Service,
	sessions *auth.Auth,
	logger *slog.Logger,
	middlewaresService *middlewares.Middlewares,
	clientInfoService *clientinfo.Service,
) Registration {
	routes := []registerRoute{
		newCreate(CreateOpts{
			Config: config, Service: service, Sessions: sessions, Logger: logger,
			Middlewares: middlewaresService, ClientInfoService: clientInfoService,
		}),
		newProfile(ProfileOpts{Service: service, Sessions: sessions}),
		newDelete(DeleteOpts{Service: service, Sessions: sessions}),
		newDeleteByKey(DeleteByKeyOpts{Service: service}),
		newServe(ServeOpts{Service: service, Logger: logger}),
	}
	for _, route := range routes {
		route.Register(api)
	}

	return Registration{}
}

type UploadedFileOutputDto struct {
	Id        string    `json:"id"`
	Name      *string   `json:"name"`
	Type      string    `json:"type"`
	Ext       string    `json:"ext"`
	Size      int64     `json:"size"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type uploadedFileWithDeleteLinkDto struct {
	UploadedFileOutputDto
	DeleteLink string `json:"delete_link"`
}

type profileOutputDto struct {
	Files []UploadedFileOutputDto `json:"files"`
	Total int                     `json:"total"`
}

type deleteOutputDto struct {
	Success bool `json:"success"`
}

func toOutput(service *uploaderservice.Service, file uploadedfile.Entity) UploadedFileOutputDto {
	return UploadedFileOutputDto{
		Id:        file.PublicID,
		Name:      file.FileName,
		Type:      file.MimeType,
		Ext:       file.Extension,
		Size:      file.SizeBytes,
		Link:      service.BuildPublicURL(file),
		CreatedAt: file.CreatedAt,
		ExpiresAt: file.ExpiresAt,
	}
}

var (
	errUploadedFileForbidden = errors.New("you don't have permission to manage this file")
)

func isExpired(file uploadedfile.Entity) bool {
	return time.Now().After(file.ExpiresAt)
}

func containsID(ids []string, id string) bool {
	for _, candidate := range ids {
		if candidate == id {
			return true
		}
	}
	return false
}
