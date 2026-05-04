package platform

import (
	"database/sql/driver"
	"fmt"
)

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

func (p *Platform) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*p = Platform(v)
	case []byte:
		*p = Platform(v)
	case nil:
		*p = ""
	default:
		return fmt.Errorf("platform: cannot scan type %T into Platform", src)
	}
	return nil
}

func (p Platform) Value() (driver.Value, error) {
	return string(p), nil
}

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
