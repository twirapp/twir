package platform

type Platform string

const (
	PlatformTwitch Platform = "twitch"
	PlatformKick   Platform = "kick"
)

func (p Platform) IsValid() bool {
	switch p {
	case PlatformTwitch, PlatformKick:
		return true
	}
	return false
}

func (p Platform) String() string { return string(p) }

func ShouldExecute(platforms []Platform, current Platform) bool {
	if len(platforms) == 0 {
		return true
	}

	for _, p := range platforms {
		if p == current {
			return true
		}
	}

	return false
}
