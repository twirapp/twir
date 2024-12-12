package resolvers

func (r *Resolver) computeBadgeUrl(fileName string) string {
	if r.config.AppEnv == "development" {
		return r.config.S3PublicUrl + "/" + r.config.S3Bucket + "/badges/" + fileName
	}

	return r.config.S3PublicUrl + "/badges/" + fileName
}
