package minio

import (
	"context"
	"fmt"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	cfg "github.com/twirapp/twir/libs/config"
)

func New(l *slog.Logger, config cfg.Config, lc *lifecycle.Lifecycle) (*minio.Client, error) {
	var creds *credentials.Credentials
	if config.AppEnv != "production" {
		creds = credentials.NewStaticV4("minio", "minio-password", "")
	} else {
		creds = credentials.NewStaticV4(config.S3AccessToken, config.S3SecretToken, "")
	}

	l.Info(
		"Creating minio client",
		slog.String("host", config.S3Host),
		slog.String("region", config.S3Region),
		slog.String("bucket", config.S3Bucket),
	)

	client, err := minio.New(
		config.S3Host,
		&minio.Options{
			Creds:  creds,
			Region: config.S3Region,
			Secure: config.AppEnv == "production",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create minio host: %w", err)
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				buckets, err := client.ListBuckets(ctx)
				if err != nil {
					return fmt.Errorf("cannot list buckets: %w", err)
				}

				bucketExists := false
				for _, bucket := range buckets {
					if bucket.Name == config.S3Bucket {
						bucketExists = true
						break
					}
				}

				if !bucketExists {
					err = client.MakeBucket(ctx, config.S3Bucket, minio.MakeBucketOptions{})
					if err != nil {
						return fmt.Errorf("cannot create bucket: %w", err)
					}
				}

				// we use cloudflare r2, which doesn't support this operation
				if config.AppEnv != "production" {
					err = client.SetBucketPolicy(
						ctx,
						config.S3Bucket,
						`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": ["s3:GetObject"],
				"Resource": [
					"arn:aws:s3:::`+config.S3Bucket+`/**"
				]
			}
		]
	}`,
					)

					if err != nil {
						return fmt.Errorf("cannot set bucket policy: %w", err)
					}

				}

				return nil
			},
			OnStop: nil,
		},
	)

	return client, nil
}

// UploaderS3Client is a dedicated S3 client for the uploader feature, separate from the shared CDN bucket client.
type UploaderS3Client struct {
	*minio.Client
	Bucket string
}

func NewUploaderS3(l *slog.Logger, config cfg.Config, lc *lifecycle.Lifecycle) (*UploaderS3Client, error) {
	host := config.UploaderS3Host
	if host == "" {
		host = config.S3Host
	}
	region := config.UploaderS3Region
	if region == "" {
		region = config.S3Region
	}
	accessToken := config.UploaderS3AccessToken
	if accessToken == "" {
		accessToken = config.S3AccessToken
	}
	secretToken := config.UploaderS3SecretToken
	if secretToken == "" {
		secretToken = config.S3SecretToken
	}

	var creds *credentials.Credentials
	if config.AppEnv != "production" {
		creds = credentials.NewStaticV4("minio", "minio-password", "")
	} else {
		creds = credentials.NewStaticV4(accessToken, secretToken, "")
	}

	l.Info(
		"Creating uploader minio client",
		slog.String("host", host),
		slog.String("region", region),
		slog.String("bucket", config.UploaderS3Bucket),
	)

	client, err := minio.New(
		host,
		&minio.Options{
			Creds:  creds,
			Region: region,
			Secure: config.AppEnv == "production",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create uploader minio host: %w", err)
	}

	uploaderClient := &UploaderS3Client{Client: client, Bucket: config.UploaderS3Bucket}
	lc.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			buckets, err := client.ListBuckets(ctx)
			if err != nil {
				return fmt.Errorf("cannot list uploader buckets: %w", err)
			}

			bucketExists := false
			for _, bucket := range buckets {
				if bucket.Name == uploaderClient.Bucket {
					bucketExists = true
					break
				}
			}
			if !bucketExists {
				if err := client.MakeBucket(ctx, uploaderClient.Bucket, minio.MakeBucketOptions{}); err != nil {
					return fmt.Errorf("cannot create uploader bucket: %w", err)
				}
			}

			if config.AppEnv != "production" {
				err := client.SetBucketPolicy(ctx, uploaderClient.Bucket, publicReadPolicy(uploaderClient.Bucket))
				if err != nil {
					return fmt.Errorf("cannot set uploader bucket policy: %w", err)
				}
			}

			return nil
		},
	})

	return uploaderClient, nil
}

func publicReadPolicy(bucket string) string {
	return `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + bucket + `/**"]
			}
		]
	}`
}
