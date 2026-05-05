package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the alert payload sent to the webhook endpoint.
type Payload struct {
	JobName   string    `json:"job_name"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Client sends alert payloads to a configured webhook URL.
type Client struct {
	URL        string
	httpClient *http.Client
}

// New creates a new webhook Client with the given URL and a sensible timeout.
func New(url string) *Client {
	return &Client{
		URL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload and posts it to the webhook URL.
// It returns an error if the request fails or the server responds with a
// non-2xx status code.
func (c *Client) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, c.URL)
	}

	return nil
}
