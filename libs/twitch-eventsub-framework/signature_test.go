package eventsub_framework

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeValidRequest() (*http.Request, []byte) {
	body := []byte("{\"test\": 123}")
	bodyReader := io.NopCloser(bytes.NewBuffer(body))

	req := httptest.NewRequest("POST", "/path/to/api", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Twitch-Eventsub-Message-Id", "286b83eb-06cb-48be-93b5-70cc1bc2a1a0")
	req.Header.Set("Twitch-Eventsub-Message-Timestamp", "2019-11-16T10:11:12.123Z")
	req.Header.Set("Twitch-Eventsub-Message-Signature", "sha256=9c496057037c816da5f1b3bfe67a6333b48868f271fbe1a282435de4f22ee71b")
	req.Header.Set("Twitch-Eventsub-Message-Retry", "0")
	req.Header.Set("Twitch-Eventsub-Message-Type", "notification")
	req.Header.Set("Twitch-Eventsub-Subscription-Type", "channel.follow")
	req.Header.Set("Twitch-Eventsub-Subscription-Version", "1")

	return req, body
}

func TestVerifyRequestSignature(t *testing.T) {
	secret := []byte("fe16e98fe09b472086054e15d0f8ac02")
	req1, body := makeValidRequest()

	valid, err :=  VerifyRequestSignature(req1, body, secret)
	if !valid || err != nil {
		t.Fatalf(`VerifyRequestSignature(req1, body, secret) = %t, %v, wanted true, nil`, valid, err)
	}

	req2, body := makeValidRequest()
	req2.Header.Del("Twitch-Eventsub-Message-Signature")
	valid, err =  VerifyRequestSignature(req2, body, secret)
	if valid || err == nil {
		t.Fatalf(`VerifyRequestSignature(req2, body, secret) = %t, %v, wanted false, missing signature header`, valid, err)
	}

	req3, body := makeValidRequest()
	req3.Header.Set("Twitch-Eventsub-Message-Signature", "abc123")
	valid, err =  VerifyRequestSignature(req3, body, secret)
	if valid || err == nil {
		t.Fatalf(`VerifyRequestSignature(req3, body, secret) = %t, %v, wanted false, invalid signature format`, valid, err)
	}
}
