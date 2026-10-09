// Package client implements the authenticated Pangolin Integration API transport.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Object = map[string]any
type Client struct {
	endpoint  *url.URL
	apiKey    string
	HTTP      *http.Client
	UserAgent string
}
type Envelope struct {
	Data    json.RawMessage `json:"data"`
	Success *bool           `json:"success"`
	Error   bool            `json:"error"`
	Status  int             `json:"status"`
}
type APIError struct {
	Status       int
	Method, Path string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Pangolin %s %s failed (status %d)", e.Method, e.Path, e.Status)
}
func IsNotFound(err error) bool  { var e *APIError; return errors.As(err, &e) && e.Status == 404 }
func NotFound(path string) error { return &APIError{Status: 404, Method: "GET", Path: path} }

func New(endpoint, apiKey string, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint")
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("endpoint must be an absolute HTTP(S) API URL without credentials, query, or fragment")
	}
	u.Path = strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(u.Path, "/v1") {
		return nil, fmt.Errorf("endpoint must include the /v1 API prefix")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	if timeout <= 0 {
		return nil, fmt.Errorf("request_timeout must be positive")
	}
	return &Client{endpoint: u, apiKey: apiKey, HTTP: &http.Client{Timeout: timeout, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}, UserAgent: "terraform-provider-pangolin"}, nil
}

// Do decodes only the envelope's data into the caller's typed destination.
// Error bodies and request/response secrets are deliberately never logged.
func (c *Client) Do(ctx context.Context, method, path string, body any, out any) error {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return fmt.Errorf("API path must start with a single slash")
	}
	rel, err := url.Parse(path)
	if err != nil || rel.IsAbs() {
		return fmt.Errorf("invalid API path")
	}
	u := *c.endpoint
	u.Path += rel.Path
	u.RawPath = c.endpoint.EscapedPath() + rel.EscapedPath()
	u.RawQuery = rel.RawQuery
	var encoded []byte
	if body != nil {
		encoded, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode API request: %w", err)
		}
	}
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(encoded))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		res, err := c.HTTP.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("Pangolin request failed: %s %s", method, rel.Path)
		}
		raw, readErr := io.ReadAll(io.LimitReader(res.Body, 16<<20))
		res.Body.Close()
		if readErr != nil {
			return fmt.Errorf("read Pangolin response")
		}
		if method == http.MethodGet && attempt < 2 && (res.StatusCode == 429 || res.StatusCode == 502 || res.StatusCode == 503 || res.StatusCode == 504) {
			delay := time.Duration(1<<attempt) * time.Second
			if n, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && n > 0 {
				delay = time.Duration(n) * time.Second
			}
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				continue
			}
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return &APIError{res.StatusCode, method, rel.Path}
		}
		if res.StatusCode == 204 && out == nil {
			return nil
		}
		var env Envelope
		if err = json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("Pangolin returned invalid JSON for %s %s", method, rel.Path)
		}
		if env.Success == nil {
			return fmt.Errorf("Pangolin response is missing success envelope for %s %s", method, rel.Path)
		}
		if !*env.Success || env.Error || env.Status >= 400 {
			status := env.Status
			if status < 400 {
				status = 500
			}
			return &APIError{status, method, rel.Path}
		}
		if out == nil {
			return nil
		}
		if len(env.Data) == 0 || bytes.Equal(env.Data, []byte("null")) {
			if method != http.MethodGet {
				return nil
			}
			return fmt.Errorf("Pangolin response is missing data for %s %s", method, rel.Path)
		}
		dec := json.NewDecoder(bytes.NewReader(env.Data))
		dec.UseNumber()
		if err = dec.Decode(out); err != nil {
			return fmt.Errorf("unexpected Pangolin response shape for %s %s", method, rel.Path)
		}
		return nil
	}
}

var parameter = regexp.MustCompile(`\{([^}]+)\}`)

func ExpandPath(template string, values map[string]any) (string, error) {
	var missing string
	p := parameter.ReplaceAllStringFunc(template, func(s string) string {
		k := s[1 : len(s)-1]
		v, ok := values[k]
		if !ok || v == nil || fmt.Sprint(v) == "" {
			missing = k
			return s
		}
		return url.PathEscape(fmt.Sprint(v))
	})
	if missing != "" {
		return "", fmt.Errorf("missing path attribute %s", missing)
	}
	return p, nil
}

// List consumes every page; missing/malformed collections are errors, never an
// empty result that could incorrectly remove a resource from Terraform state.
func (c *Client) List(ctx context.Context, path, key, pagination string) ([]Object, error) {
	var all []Object
	offset, page := 0, 1
	// A small number of pinned upstream handlers report an unfiltered total.
	// Their definitions explicitly opt into page-length termination.
	ignoreTotal := strings.HasSuffix(pagination, "-short")
	pagination = strings.TrimSuffix(pagination, "-short")
	for requests := 0; requests < 10000; requests++ {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		if pagination == "page" {
			q.Set("page", strconv.Itoa(page))
			q.Set("pageSize", "100")
		} else if pagination == "offset" {
			q.Set("limit", "100")
			q.Set("offset", strconv.Itoa(offset))
		}
		u.RawQuery = q.Encode()
		var data Object
		if err = c.Do(ctx, "GET", u.String(), nil, &data); err != nil {
			return nil, err
		}
		raw, ok := data[key].([]any)
		if !ok {
			return nil, fmt.Errorf("Pangolin list response is missing %s array", key)
		}
		for _, v := range raw {
			o, ok := v.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid object in %s array", key)
			}
			all = append(all, o)
		}
		p, _ := data["pagination"].(map[string]any)
		total, totalOK := integer(p["total"])
		if ignoreTotal {
			totalOK = false
		}
		if len(raw) == 0 && totalOK && len(all) < total {
			return nil, fmt.Errorf("Pangolin returned an incomplete %s listing", key)
		}
		if pagination == "" || len(raw) == 0 || (totalOK && len(all) >= total) {
			return all, nil
		}
		if !totalOK && len(raw) < 100 {
			return all, nil
		}
		offset += len(raw)
		page++
	}
	return nil, fmt.Errorf("Pangolin pagination exceeded safety limit")
}
func integer(v any) (int, bool) {
	if v == nil {
		return 0, false
	}
	n, e := strconv.Atoi(fmt.Sprint(v))
	return n, e == nil
}
