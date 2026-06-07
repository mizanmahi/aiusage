package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mizanmahi/aiusage/types"
)

const minCLIVersionHeader = "X-Aiusage-Min-CLI-Version"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func Send(serverURL, apiKey, clientVersion string, events []types.UsageEvent) (*types.PushResponse, error) {
	return NewClient().Send(serverURL, apiKey, clientVersion, events)
}

func (c *Client) Send(serverURL, apiKey, clientVersion string, events []types.UsageEvent) (*types.PushResponse, error) {
	body, err := json.Marshal(types.PushPayload{Events: events})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(serverURL, "/")+"/ingest", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	if clientVersion != "" {
		req.Header.Set("X-Aiusage-CLI-Version", clientVersion)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkCompatibility(resp, clientVersion); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	var result types.APIResponse[types.PushResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func checkCompatibility(resp *http.Response, clientVersion string) error {
	minVersion := resp.Header.Get(minCLIVersionHeader)
	if minVersion == "" || clientVersion == "" {
		return nil
	}
	if compareSemver(clientVersion, minVersion) >= 0 {
		return nil
	}

	return fmt.Errorf("server requires aiusage CLI >= %s, current version is %s", minVersion, clientVersion)
}

func compareSemver(left, right string) int {
	leftParts := parseSemver(left)
	rightParts := parseSemver(right)

	for i := range leftParts {
		if leftParts[i] > rightParts[i] {
			return 1
		}
		if leftParts[i] < rightParts[i] {
			return -1
		}
	}

	return 0
}

func parseSemver(version string) [3]int {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	parts := strings.Split(version, ".")

	var parsed [3]int
	for i := 0; i < len(parts) && i < len(parsed); i++ {
		for _, digit := range parts[i] {
			if digit < '0' || digit > '9' {
				break
			}
			parsed[i] = parsed[i]*10 + int(digit-'0')
		}
	}

	return parsed
}
