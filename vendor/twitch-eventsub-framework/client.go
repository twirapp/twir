package eventsub_framework

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	"net/http"
	"net/url"
	"time"
)

const (
	EventSubSubscriptionsEndpoint = "https://api.twitch.tv/helix/eventsub/subscriptions"

	pageSize = "100"
)

type Credentials interface {
	ClientID() (string, error)
	AppToken() (string, error)
}

type SubRequest struct {
	// The type of event being subscribed to.
	Type string
	// The parameters under which the event will be fired.
	Condition interface{}
	// The Webhook HTTP callback address.
	Callback string
	// The HMAC secret used to verify the event data.
	Secret string
}

type Status string

const (
	StatusAny                  Status = ""
	StatusEnabled              Status = "enabled"
	StatusVerificationPending  Status = "webhook_callback_verification_pending"
	StatusVerificationFailed   Status = "webhook_callback_verification_failed"
	StatusFailuresExceeded     Status = "notification_failures_exceeded"
	StatusAuthorizationRevoked Status = "authorization_revoked"
	StatusModeratorRemoved     Status = "moderator_removed"
	StatusUserRemoved          Status = "user_removed"
	StatusVersionRemoved       Status = "version_removed"
)

// TwitchError describes an error from the Twitch API.
//
// For example:
//  {
//    "error": "Unauthorized",
//    "status": 401,
//    "message": "Invalid OAuth token"
//  }
type TwitchError struct {
	ErrorText string `json:"error"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

func (t *TwitchError) Error() string {
	if t.Message != "" {
		return fmt.Sprintf("%d %s: %s", t.Status, t.ErrorText, t.Message)
	} else {
		return fmt.Sprintf("%d %s", t.Status, t.ErrorText)
	}
}

type SubClient struct {
	httpClient  *http.Client
	credentials Credentials
}

// NewSubClient creates a new SubClient with the given Credentials provider.
func NewSubClient(credentials Credentials) *SubClient {
	return &SubClient{
		httpClient: &http.Client{
			Timeout: time.Second * 3,
		},
		credentials: credentials,
	}
}

// NewSubClientHTTP creates a new SubClient with the given Credentials provider
// and http.Client instance.
func NewSubClientHTTP(credentials Credentials, client *http.Client) *SubClient {
	return &SubClient{
		httpClient:  client,
		credentials: credentials,
	}
}

// Performs a given http.Request while adding the Client-ID and Authorization
// headers to the request.
//
// If the returned error is non-nil, the caller must  close the returned
// response body. The returned response is guaranteed to have a 2xx status code.
func (s *SubClient) do(req *http.Request) (*http.Response, error) {
	clientID, err := s.credentials.ClientID()
	if err != nil {
		return nil, fmt.Errorf("get client id: %w", err)
	}

	appToken, err := s.credentials.AppToken()
	if err != nil {
		return nil, fmt.Errorf("get app token: %w", err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+appToken)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		defer res.Body.Close()
		var twitchErr TwitchError
		if err := json.NewDecoder(res.Body).Decode(&twitchErr); err != nil {
			return nil, fmt.Errorf("process %d twitch api status: %w", res.StatusCode, err)
		}
		return nil, &twitchErr
	}

	return res, nil
}

// Subscribe creates a new Webhook subscription.
func (s *SubClient) Subscribe(ctx context.Context, srq *SubRequest) (*esb.RequestStatus, error) {
	reqJSON := esb.Request{
		Type:      srq.Type,
		Version:   "1",
		Condition: srq.Condition,
		Transport: esb.Transport{
			Method:   "webhook",
			Callback: srq.Callback,
			Secret:   srq.Secret,
		},
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqJSON)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", EventSubSubscriptionsEndpoint, buf)
	if err != nil {
		return nil, err
	}
	res, err := s.do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var statusResponse esb.RequestStatus
	if err := json.NewDecoder(res.Body).Decode(&statusResponse); err != nil {
		return nil, err
	}

	return &statusResponse, nil
}

// Unsubscribe deletes a Webhook subscription by the subscription's ID.
func (s *SubClient) Unsubscribe(ctx context.Context, subscriptionID string) error {
	u, err := url.Parse(EventSubSubscriptionsEndpoint)
	if err != nil {
		return fmt.Errorf("unsubscribe: parse EventSubSubscriptionsEndpoint url: %w", err)
	}

	q := u.Query()
	q.Set("id", subscriptionID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	res, err := s.do(req)
	if err != nil {
		return err
	}

	// Performing the HTTP request did not return an error, so the status code
	// must have been a 2xx. No need to worry about reading the (possible) body.
	_ = res.Body.Close()
	return nil
}

// GetSubscriptions returns all EventSub subscriptions.
// If statusFilter != StatusAny, it will apply the filter to the query.
func (s *SubClient) GetSubscriptions(ctx context.Context, statusFilter Status) (*esb.RequestStatus, error) {
	firstRes, err := s.getSubscriptions(ctx, statusFilter, "")
	if err != nil {
		return nil, err
	}

	if firstRes.Pagination == nil || firstRes.Pagination.Cursor == "" {
		// No pagination was specified.
		return firstRes, nil
	}

	cursor := firstRes.Pagination.Cursor

	// arbitrary number over 100, the maximum number of pages
	for i := 1; i < 105; i++ {
		nextRes, err := s.getSubscriptions(ctx, statusFilter, cursor)
		if err != nil {
			return nil, err
		}

		// Combine data from each page into firstReq
		firstRes.Data = append(firstRes.Data, nextRes.Data...)

		if nextRes.Pagination == nil || nextRes.Pagination.Cursor == "" {
			// Done with all the pages
			return firstRes, nil
		} else {
			cursor = nextRes.Pagination.Cursor
		}
		i++
	}

	return nil, fmt.Errorf("caught in loop while following pagination")
}

// Get the subscriptions with a specific pagination cursor
func (s *SubClient) getSubscriptions(ctx context.Context, statusFilter Status, cursor string) (*esb.RequestStatus, error) {
	// First, construct the request url with the proper query parameters.
	u, err := url.Parse(EventSubSubscriptionsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("get subscriptions: parse EventSubSubscriptionsEndpoint url: %w", err)
	}

	q := u.Query()
	q.Set("first", pageSize)
	if statusFilter != StatusAny {
		q.Set("status", string(statusFilter))
	}
	if cursor != "" {
		q.Set("after", cursor)
	}
	u.RawQuery = q.Encode()
	reqUrl := u.String()

	// Now, actually send the request.
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var subscriptionsResponse esb.RequestStatus
	if err := json.NewDecoder(res.Body).Decode(&subscriptionsResponse); err != nil {
		return nil, err
	}

	return &subscriptionsResponse, nil
}
