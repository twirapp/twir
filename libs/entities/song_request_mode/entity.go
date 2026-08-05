package song_request_mode

import (
	"database/sql/driver"
	"fmt"
	"slices"
)

type Mode string

const (
	ModeYouTube Mode = "YOUTUBE"
	ModeSpotify Mode = "SPOTIFY"
)

var allModes = []Mode{ModeYouTube, ModeSpotify}

func (m Mode) IsValid() bool {
	return slices.Contains(allModes, m)
}

func (m Mode) String() string { return string(m) }

func (m *Mode) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*m = Mode(v)
	case []byte:
		*m = Mode(v)
	case nil:
		*m = ""
	default:
		return fmt.Errorf("song request mode: cannot scan type %T into Mode", src)
	}
	return nil
}

func (m Mode) Value() (driver.Value, error) {
	return string(m), nil
}

var ErrInvalidMode = fmt.Errorf("invalid song request mode")
