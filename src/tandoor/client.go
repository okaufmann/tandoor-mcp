package tandoor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/dto"
)


type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL, apiToken string, logHTTPBody bool) *Client {
	baseURL = strings.TrimSuffix(baseURL, "/")

	authHeader := "Bearer " + apiToken

	var transport http.RoundTripper = &authTransport{
		token: authHeader,
		next:  http.DefaultTransport,
	}

	if logHTTPBody {
		transport = &loggingTransport{
			next: transport,
		}
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
	}
}

// checkResponse parses any error responses and acts as a dedicated error sink
func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	return &dto.APIError{
		StatusCode: resp.StatusCode,
		Message:    string(bodyBytes),
	}
}

// Request executes an HTTP request using the Client and parses the JSON response into type T.
func Request[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) (*T, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path
	u.RawQuery = query.Encode()

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var response T
	if resp.StatusCode == http.StatusNoContent {
		return &response, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		if err == io.EOF {
			return &response, nil
		}
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &response, nil
}
