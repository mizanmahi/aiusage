package push

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mizanmahi/aiusage/types"
)

func TestSendPostsEventsWithBearerAuth(t *testing.T) {
	var gotAuth string
	var gotVersion string
	var gotPayload types.PushPayload

	client := &Client{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/ingest" {
			t.Fatalf("path = %q, want /ingest", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}

		gotAuth = r.Header.Get("Authorization")
		gotVersion = r.Header.Get("X-Aiusage-CLI-Version")
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}

		body := mustJSON(t, types.APIResponse[types.PushResponse]{
			Data: types.PushResponse{Accepted: 1, Skipped: 0, Message: "ok"},
		})
		return jsonResponse(http.StatusOK, body, nil), nil
	})}}

	result, err := client.Send("http://example.test/", "ak_secret", "1.2.3", []types.UsageEvent{{SessionID: "session-1"}})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if gotAuth != "Bearer ak_secret" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
	if gotVersion != "1.2.3" {
		t.Fatalf("X-Aiusage-CLI-Version = %q, want 1.2.3", gotVersion)
	}
	if len(gotPayload.Events) != 1 || gotPayload.Events[0].SessionID != "session-1" {
		t.Fatalf("payload events = %+v, want session-1", gotPayload.Events)
	}
	if result.Accepted != 1 || result.Message != "ok" {
		t.Fatalf("result = %+v, want accepted ok", result)
	}
}

func TestSendReturnsErrorForNonOKStatus(t *testing.T) {
	client := &Client{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"error":"nope"}`, nil), nil
	})}}

	_, err := client.Send("http://example.test", "ak_secret", "1.2.3", nil)
	if err == nil {
		t.Fatal("Send() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "server returned 401") {
		t.Fatalf("Send() error = %q, want status code", err.Error())
	}
}

func TestSendChecksMinimumCLIVersion(t *testing.T) {
	headers := http.Header{}
	headers.Set(minCLIVersionHeader, "2.0.0")
	client := &Client{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := mustJSON(t, types.APIResponse[types.PushResponse]{
			Data: types.PushResponse{Accepted: 1, Message: "ok"},
		})
		return jsonResponse(http.StatusOK, body, headers), nil
	})}}

	_, err := client.Send("http://example.test", "ak_secret", "1.5.0", nil)
	if err == nil {
		t.Fatal("Send() error = nil, want compatibility error")
	}
	if !strings.Contains(err.Error(), "server requires aiusage CLI >= 2.0.0") {
		t.Fatalf("Send() error = %q, want compatibility message", err.Error())
	}
}

func TestSendAllowsSatisfiedMinimumCLIVersion(t *testing.T) {
	headers := http.Header{}
	headers.Set(minCLIVersionHeader, "1.2.0")
	client := &Client{httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := mustJSON(t, types.APIResponse[types.PushResponse]{
			Data: types.PushResponse{Accepted: 1, Message: "ok"},
		})
		return jsonResponse(http.StatusOK, body, headers), nil
	})}}

	result, err := client.Send("http://example.test", "ak_secret", "1.2.3", nil)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if result.Accepted != 1 {
		t.Fatalf("Accepted = %d, want 1", result.Accepted)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func jsonResponse(status int, body string, headers http.Header) *http.Response {
	if headers == nil {
		headers = http.Header{}
	}
	headers.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: status,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	return string(data)
}

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{name: "equal", left: "1.2.3", right: "1.2.3", want: 0},
		{name: "left greater", left: "1.3.0", right: "1.2.9", want: 1},
		{name: "left lower", left: "1.2.0", right: "1.2.1", want: -1},
		{name: "v prefix", left: "v2.0.0", right: "1.9.9", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareSemver(tt.left, tt.right)
			if got != tt.want {
				t.Fatalf("compareSemver(%q, %q) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
