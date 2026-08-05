package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestSpotify_GetDevices_returns_devices(t *testing.T) {
	defaultClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("method = %s, want GET", req.Method)
			}
			if req.URL.Path != "/v1/me/player/devices" {
				t.Fatalf("path = %s, want /v1/me/player/devices", req.URL.Path)
			}
			if auth := req.Header.Get("Authorization"); auth != "Bearer access-token" {
				t.Fatalf("Authorization = %q, want bearer token", auth)
			}

			body := `{"devices":[{"id":"device-1","name":"Desktop","is_active":true,"is_restricted":false,"is_private_session":false},{"id":"device-2","name":"Phone","is_active":false,"is_restricted":true,"is_private_session":true}]}`
			return spotifyTestResponse(req, http.StatusOK, body, nil), nil
		}),
	}
	t.Cleanup(func() {
		http.DefaultClient = defaultClient
	})

	devices, err := NewStatic("access-token", nil).GetDevices(context.Background())
	if err != nil {
		t.Fatalf("GetDevices() error = %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("len(devices) = %d, want 2", len(devices))
	}
	if devices[0].ID != "device-1" || !devices[0].IsActive || devices[0].IsRestricted || devices[0].IsPrivateSession {
		t.Fatalf("device[0] = %#v, want parsed fields", devices[0])
	}
	if devices[1].ID != "device-2" || devices[1].IsActive || !devices[1].IsRestricted || !devices[1].IsPrivateSession {
		t.Fatalf("device[1] = %#v, want parsed fields", devices[1])
	}
}
