package minio

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	cfg "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

func New(l *slog.Logger, config cfg.Config, lc fx.Lifecycle) (*minio.Client, error) {
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
		fx.Hook{
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
