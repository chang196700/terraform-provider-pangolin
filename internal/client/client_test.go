package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestEnvelopeAndAuthentication(t *testing.T) {
	for _, tc := range []struct {
		name          string
		http, status  int
		success       bool
		expectedError bool
	}{{"created", 201, 201, true, false}, {"business error", 200, 403, false, true}, {"http error", 403, 403, false, true}} {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/example" || r.Header.Get("Authorization") != "Bearer secret" {
					t.Errorf("incorrect authenticated request")
				}
				w.WriteHeader(tc.http)
				json.NewEncoder(w).Encode(map[string]any{"success": tc.success, "error": !tc.success, "status": tc.status, "data": Object{"id": 9007199254740993}, "message": "secret must not leak"})
			}))
			defer s.Close()
			c, _ := New(s.URL+"/v1/", "secret", time.Second)
			var o Object
			err := c.Do(context.Background(), "PUT", "/example", Object{"enabled": false}, &o)
			if (err != nil) != tc.expectedError {
				t.Fatalf("unexpected error %v", err)
			}
			if err != nil && strings.Contains(err.Error(), "secret") {
				t.Fatal("sensitive response leaked")
			}
			if err == nil && fmt.Sprint(o["id"]) != "9007199254740993" {
				t.Fatal("numeric precision lost")
			}
		})
	}
}
func TestCreateNeverRetries(t *testing.T) {
	var n atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.Add(1); w.WriteHeader(503) }))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	if c.Do(context.Background(), "PUT", "/org", Object{}, nil) == nil {
		t.Fatal("expected error")
	}
	if n.Load() != 1 {
		t.Fatal("non-idempotent create retried")
	}
}
func TestReadRetryCancellation(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Retry-After", "10"); w.WriteHeader(429) }))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := time.Now()
	if c.Do(ctx, "GET", "/orgs", nil, &Object{}) == nil {
		t.Fatal("expected cancellation")
	}
	if time.Since(started) > time.Second {
		t.Fatal("retry ignored cancellation")
	}
}
func TestPaginationBothProtocols(t *testing.T) {
	for _, kind := range []string{"offset", "page"} {
		t.Run(kind, func(t *testing.T) {
			var n int
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n++
				expected := "1"
				key := "page"
				if kind == "offset" {
					expected = "0"
					key = "offset"
				}
				if n == 2 {
					expected = "2"
					if kind == "offset" {
						expected = "100"
					}
				}
				if r.URL.Query().Get(key) != expected {
					t.Errorf("incorrect page query: %s", r.URL.RawQuery)
				}
				rows := []Object{}
				if n == 1 {
					for i := 0; i < 100; i++ {
						rows = append(rows, Object{"id": i})
					}
				} else {
					rows = append(rows, Object{"id": 999})
				}
				json.NewEncoder(w).Encode(map[string]any{"success": true, "status": 200, "data": Object{"items": rows, "pagination": Object{"total": 101}}})
			}))
			defer s.Close()
			c, _ := New(s.URL+"/v1", "key", time.Second)
			rows, err := c.List(context.Background(), "/items", "items", kind)
			if err != nil || len(rows) != 101 || fmt.Sprint(rows[100]["id"]) != "999" || n != 2 {
				t.Fatalf("pagination failed: %d %v %d", len(rows), err, n)
			}
		})
	}
}
func TestMalformedListIsNotEmpty(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{"success":true,"status":200,"data":{}}`) }))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	_, err := c.List(context.Background(), "/items", "items", "offset")
	if err == nil || IsNotFound(err) {
		t.Fatal("malformed response mistaken for an absent resource")
	}
}

func TestIncompletePaginationIsNotAbsent(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"status":200,"data":{"items":[],"pagination":{"total":10}}}`)
	}))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	_, err := c.List(context.Background(), "/items", "items", "offset")
	if err == nil || IsNotFound(err) {
		t.Fatalf("incomplete list must not become absence: %v", err)
	}
}
func TestNoCredentialRedirect(t *testing.T) {
	var received bool
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { received = true }))
	defer destination.Close()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, destination.URL, 302) }))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	err := c.Do(context.Background(), "GET", "/items", nil, &Object{})
	if err == nil || received {
		t.Fatal("authenticated redirect was followed")
	}
}
func TestNullDataMutation(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"status":200,"data":null}`)
	}))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	var o Object
	if err := c.Do(context.Background(), "POST", "/item", Object{}, &o); err != nil {
		t.Fatal(err)
	}
	if err := c.Do(context.Background(), "GET", "/item", nil, &o); err == nil {
		t.Fatal("null GET response accepted")
	}
}

func TestPathParametersStayOneSegment(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/v1/org/a%2Fb" {
			t.Errorf("escaped ID changed: %s", r.URL.EscapedPath())
		}
		fmt.Fprint(w, `{"success":true,"status":200,"data":{"id":"a/b"}}`)
	}))
	defer s.Close()
	c, _ := New(s.URL+"/v1", "key", time.Second)
	p, err := ExpandPath("/org/{id}", Object{"id": "a/b"})
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Do(context.Background(), "GET", p, nil, &Object{}); err != nil {
		t.Fatal(err)
	}
}
