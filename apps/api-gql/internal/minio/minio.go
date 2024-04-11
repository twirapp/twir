package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

func New(config cfg.Config, lc fx.Lifecycle) (*minio.Client, error) {
	var creds *credentials.Credentials
	if config.AppEnv != "production" {
		creds = credentials.NewStaticV4("minio", "minio-password", "")
	} else {
		creds = credentials.NewStaticV4(config.S3AccessToken, config.S3SecretToken, "")
	}

	client, err := minio.New(
		config.S3Host,
		&minio.Options{
			Creds:  creds,
			Region: config.S3Region,
			Secure: config.AppEnv == "production",
		},
	)
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				buckets, err := client.ListBuckets(context.TODO())
				if err != nil {
					return err
				}

				bucketExists := false
				for _, bucket := range buckets {
					if bucket.Name == config.S3Bucket {
						bucketExists = true
						break
					}
				}

				if !bucketExists {
					err = client.MakeBucket(context.TODO(), config.S3Bucket, minio.MakeBucketOptions{})
					if err != nil {
						return err
					}
				}

				return client.SetBucketPolicy(
					context.TODO(),
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
			},
			OnStop: nil,
		},
	)

	return client, nil
}
