package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type ChannelIntegrationValorant struct {
	ID            uuid.UUID
	Enabled       bool
	AccessToken   *string
	RefreshToken  *string
	ApiKey        *string
	Data          *Data
	ChannelID     string
	IntegrationID string
}

var Nil = ChannelIntegrationValorant{}

type Data struct {
	// username#tag
	UserName             *string `json:"username,omitempty"`
	ValorantActiveRegion *string `json:"valorantActiveRegion,omitempty"`
	ValorantPuuid        *string `json:"valorantPuuid,omitempty"`
}

func (u *Data) Value() (driver.Value, error) {
	if u == nil {
		return "{}", nil
	}

	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}

// Scan implements the sql.Scanner interface for Data
func (u *Data) Scan(src interface{}) error {
	// If the column is NULL, src will be nil.
	if src == nil {
		// Do nothing, the pointer is already nil or will be.
		return nil
	}

	switch sourceType := src.(type) {
	case string:
		return json.Unmarshal([]byte(sourceType), u)
	case []byte:
		return json.Unmarshal(sourceType, u)
	}

	return fmt.Errorf("unsupported Scan, storing %T into Data", src)
}
