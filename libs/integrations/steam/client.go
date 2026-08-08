package steam

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

const (
	loginURL        = "https://steamcommunity.com/openid/login"
	claimedIDPrefix = "https://steamcommunity.com/openid/id/"

	maxResponseBodyBytes = 1 << 20
)

var (
	ErrAssertionInvalid = errors.New("steam openid assertion is not valid")
	ErrNoClaimedID      = errors.New("steam openid assertion has no claimed id")
)

type PlayerSummary struct {
	SteamID     string
	PersonaName string
	AvatarFull  string
	ProfileURL  string
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c.apiKey != ""
}

// BuildAuthURL returns the Steam OpenID 2.0 login URL. Steam redirects back to
// returnTo with the signed assertion in the query string.
func BuildAuthURL(returnTo string) (string, error) {
	parsed, err := url.Parse(returnTo)
	if err != nil {
		return "", fmt.Errorf("steam: parse return url: %w", err)
	}
	realm := parsed.Scheme + "://" + parsed.Host + "/"

	params := url.Values{
		"openid.ns":         {"http://specs.openid.net/auth/2.0"},
		"openid.mode":       {"checkid_setup"},
		"openid.return_to":  {returnTo},
		"openid.realm":      {realm},
		"openid.identity":   {"http://specs.openid.net/auth/2.0/identifier_select"},
		"openid.claimed_id": {"http://specs.openid.net/auth/2.0/identifier_select"},
	}

	return loginURL + "?" + params.Encode(), nil
}

// VerifyAssertion validates openid.* response params (as received on the
// return_to URL) against Steam and returns the SteamID64 from the claimed id.
func (c *Client) VerifyAssertion(ctx context.Context, query url.Values) (string, error) {
	mode := query.Get("openid.mode")
	if mode != "id_res" {
		return "", fmt.Errorf("unexpected openid.mode %q", mode)
	}

	claimedID := query.Get("openid.claimed_id")
	if claimedID == "" {
		return "", ErrNoClaimedID
	}

	verification := make(url.Values, len(query)+1)
	for key, values := range query {
		verification[key] = values
	}
	verification.Set("openid.mode", "check_authentication")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		loginURL,
		strings.NewReader(verification.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("steam: create check_authentication request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("steam: check_authentication request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return "", fmt.Errorf("steam: read check_authentication response: %w", err)
	}

	if !strings.Contains(string(body), "is_valid:true") {
		return "", ErrAssertionInvalid
	}

	steamID64 := strings.TrimPrefix(claimedID, claimedIDPrefix)
	if steamID64 == claimedID {
		return "", fmt.Errorf("steam: unexpected claimed id format %q", claimedID)
	}

	return steamID64, nil
}

var steamID64Pattern = regexp.MustCompile(`^\d{17}$`)

func (c *Client) GetPlayerSummaries(ctx context.Context, steamID64s []string) ([]PlayerSummary, error) {
	if !c.Enabled() {
		return nil, nil
	}
	for _, id := range steamID64s {
		if !steamID64Pattern.MatchString(id) {
			return nil, fmt.Errorf("steam: invalid steam id %q", id)
		}
	}

	query := url.Values{
		"key":      {c.apiKey},
		"steamids": {strings.Join(steamID64s, ",")},
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("steam: create player summaries request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("steam: player summaries request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("steam: read player summaries response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam: player summaries status %d", resp.StatusCode)
	}

	var parsed struct {
		Response struct {
			Players []struct {
				SteamID     string `json:"steamid"`
				PersonaName string `json:"personaname"`
				AvatarFull  string `json:"avatarfull"`
				ProfileURL  string `json:"profileurl"`
			} `json:"players"`
		} `json:"response"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("steam: parse player summaries: %w", err)
	}

	summaries := make([]PlayerSummary, 0, len(parsed.Response.Players))
	for _, player := range parsed.Response.Players {
		summaries = append(summaries, PlayerSummary{
			SteamID:     player.SteamID,
			PersonaName: player.PersonaName,
			AvatarFull:  player.AvatarFull,
			ProfileURL:  player.ProfileURL,
		})
	}

	return summaries, nil
}
