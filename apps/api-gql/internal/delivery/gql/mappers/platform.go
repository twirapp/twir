package mappers

import "github.com/twirapp/twir/libs/entities/platform"

func PlatformsToStrings(platforms []platform.Platform) []string {
	if platforms == nil {
		return []string{}
	}
	result := make([]string, len(platforms))
	for i, p := range platforms {
		result[i] = p.String()
	}
	return result
}

func StringsToPlatforms(ss []string) []platform.Platform {
	if ss == nil {
		return []platform.Platform{}
	}
	result := make([]platform.Platform, len(ss))
	for i, s := range ss {
		result[i] = platform.Platform(s)
	}
	return result
}
