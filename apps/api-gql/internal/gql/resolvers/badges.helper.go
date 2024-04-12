package resolvers

func (r *Resolver) computeBadgeUrl(id string) string {
	if r.config.AppEnv == "development" {
		return r.config.S3PublicUrl + "/" + r.config.S3Bucket + "/badges/" + id
	}

	return r.config.S3Host + "/badges/" + id
}
