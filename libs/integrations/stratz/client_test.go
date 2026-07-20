package stratz

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWinProbability_RejectsOversizedResponseBody(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(
				w,
				`{"data":`+
					strings.Repeat(" ", 1<<20)+
					`{"live":{"match":{"liveWinRateValues":[{"time":1,"winRate":0.625}]}}}}`,
			)
		}),
	)
	t.Cleanup(server.Close)

	_, err := New("token", WithBaseURL(server.URL)).WinProbability(context.Background(), 1)
	if err == nil {
		t.Fatal("expected oversized response body error")
	}
	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Errorf("expected response body limit error, got %v", err)
	}
}

func TestWinProbability_TruncatesNonSuccessResponsePreview(t *testing.T) {
	body := strings.Repeat("x", 8<<10)
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, body, http.StatusInternalServerError)
		}),
	)
	t.Cleanup(server.Close)

	_, err := New("token", WithBaseURL(server.URL)).WinProbability(context.Background(), 1)
	if err == nil {
		t.Fatal("expected unexpected status error")
	}
	if !strings.HasSuffix(err.Error(), "...") {
		t.Errorf("expected truncated response preview, got %q", err)
	}
	if strings.Contains(err.Error(), body) {
		t.Error("error included the full response body")
	}
}

func TestWinProbability_NormalizesWinRateScale(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     float64
		wantErr  bool
	}{
		{
			name: "fractional series",
			response: `{"data":{"live":{"match":{"liveWinRateValues":[
				{"time":100,"winRate":0.41},{"time":200,"winRate":0.625}
			]}}}}`,
			want: 0.625,
		},
		{
			name: "percentage series",
			response: `{"data":{"live":{"match":{"liveWinRateValues":[
				{"time":100,"winRate":41.2},{"time":200,"winRate":62.5}
			]}}}}`,
			want: 0.625,
		},
		{
			name: "percentage series with fractional final value",
			response: `{"data":{"live":{"match":{"liveWinRateValues":[
				{"time":100,"winRate":41.2},{"time":200,"winRate":0.625}
			]}}}}`,
			want: 0.00625,
		},
		{
			name: "out of range final value",
			response: `{"data":{"live":{"match":{"liveWinRateValues":[
				{"time":100,"winRate":101}
			]}}}}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, tt.response)
				}),
			)
			t.Cleanup(server.Close)

			probability, err := New("token", WithBaseURL(server.URL)).WinProbability(context.Background(), 1)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected out of range error")
				}
				if !strings.Contains(err.Error(), "out of range") {
					t.Errorf("expected out of range error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("WinProbability returned error: %v", err)
			}
			if probability != tt.want {
				t.Errorf("expected %v, got %v", tt.want, probability)
			}
		})
	}
}
