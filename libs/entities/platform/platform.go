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
