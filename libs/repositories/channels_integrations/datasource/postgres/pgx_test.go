package postgres

import (
	"reflect"
	"testing"

	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	"github.com/twirapp/twir/libs/repositories/channels_integrations/model"
)

func TestMarshalDataUsesCanonicalProfileKeys(t *testing.T) {
	t.Parallel()

	userName := "streamer"
	avatar := "https://avatar.test/profile.png"
	data, err := marshalData(&model.Data{UserName: &userName, Avatar: &avatar})
	if err != nil {
		t.Fatalf("marshalData() error = %v", err)
	}
	if got, want := string(data), `{"avatar":"https://avatar.test/profile.png","username":"streamer"}`; got != want {
		t.Fatalf("marshalData() = %q, want %q", got, want)
	}
}

func TestBuildUpdateQuerySetsTokens(t *testing.T) {
	t.Parallel()

	accessToken := "new-access"
	refreshToken := "new-refresh"
	query, args, err := buildUpdateQuery("integration-id", channelsintegrations.UpdateInput{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	})
	if err != nil {
		t.Fatalf("buildUpdateQuery() error = %v", err)
	}
	if want := `UPDATE channels_integrations SET "accessToken" = $1, "refreshToken" = $2 WHERE id = $3`; query != want {
		t.Fatalf("buildUpdateQuery() query = %q, want %q", query, want)
	}
	if want := []any{"new-access", "new-refresh", "integration-id"}; !reflect.DeepEqual(args, want) {
		t.Fatalf("buildUpdateQuery() args = %#v, want %#v", args, want)
	}
}

func TestBuildUpdateQueryClearsTokens(t *testing.T) {
	t.Parallel()

	query, args, err := buildUpdateQuery("integration-id", channelsintegrations.UpdateInput{
		ClearAccessToken:  true,
		ClearRefreshToken: true,
	})
	if err != nil {
		t.Fatalf("buildUpdateQuery() error = %v", err)
	}
	if want := `UPDATE channels_integrations SET "accessToken" = NULL, "refreshToken" = NULL WHERE id = $1`; query != want {
		t.Fatalf("buildUpdateQuery() query = %q, want %q", query, want)
	}
	if want := []any{"integration-id"}; !reflect.DeepEqual(args, want) {
		t.Fatalf("buildUpdateQuery() args = %#v, want %#v", args, want)
	}
}

func TestBuildUpdateQueryPreservesTokensForUnrelatedUpdates(t *testing.T) {
	t.Parallel()

	enabled := false
	query, args, err := buildUpdateQuery("integration-id", channelsintegrations.UpdateInput{Enabled: &enabled})
	if err != nil {
		t.Fatalf("buildUpdateQuery() error = %v", err)
	}
	if want := `UPDATE channels_integrations SET enabled = $1 WHERE id = $2`; query != want {
		t.Fatalf("buildUpdateQuery() query = %q, want %q", query, want)
	}
	if want := []any{false, "integration-id"}; !reflect.DeepEqual(args, want) {
		t.Fatalf("buildUpdateQuery() args = %#v, want %#v", args, want)
	}
}

func TestBuildUpdateQueryRejectsSettingAndClearingSameToken(t *testing.T) {
	t.Parallel()

	accessToken := "new-access"
	refreshToken := "new-refresh"
	tests := []struct {
		name  string
		input channelsintegrations.UpdateInput
	}{
		{
			name: "access token",
			input: channelsintegrations.UpdateInput{
				AccessToken: &accessToken, ClearAccessToken: true,
			},
		},
		{
			name: "refresh token",
			input: channelsintegrations.UpdateInput{
				RefreshToken: &refreshToken, ClearRefreshToken: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if _, _, err := buildUpdateQuery("integration-id", test.input); err == nil {
				t.Fatal("buildUpdateQuery() error = nil, want set/clear conflict")
			}
		})
	}
}
