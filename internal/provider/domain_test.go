package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func TestDomainLifecycleUsesOfficialRoutes(t *testing.T) {
	var updateBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/v1/org/org-1/domain":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["type"] != "wildcard" || body["baseDomain"] != "example.com" || body["preferWildcardCert"] != false {
				t.Fatalf("unexpected create body: %#v", body)
			}
			_, _ = w.Write([]byte(`{"success":true,"error":false,"status":201,"data":{"domainId":"dom-1","certResolver":null,"preferWildcardCert":false}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/org/org-1/domain/dom-1":
			_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"domainId":"dom-1","baseDomain":"example.com","type":"wildcard","verified":true,"failed":false,"tries":0,"configManaged":false,"certResolver":null,"customCertResolver":"letsencrypt","preferWildcardCert":false,"errorMessage":null,"lastCheckedAt":100}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/org/org-1/domain/dom-1":
			if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
				t.Fatal(err)
			}
			_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"domainId":"dom-1","certResolver":"custom","preferWildcardCert":false}}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/org/org-1/domain/dom-1":
			_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"success":true}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := domainDefinition()
	v := Values{"org_id": "org-1", "type": "wildcard", "base_domain": "example.com", "prefer_wildcard_cert": false}
	body, err := Body(d, v, true)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Call(context.Background(), c, d.Create, v, body)
	if err != nil {
		t.Fatal(err)
	}
	if Text(o["domainId"]) != "dom-1" {
		t.Fatalf("unexpected create response: %#v", o)
	}
	v["id"] = "dom-1"
	if err := d.UpdateFunc(context.Background(), c, Values{"id": "dom-1", "org_id": "org-1", "type": "wildcard", "cert_resolver": "custom", "prefer_wildcard_cert": false}, Values{"id": "dom-1", "org_id": "org-1", "type": "wildcard", "cert_resolver": nil, "prefer_wildcard_cert": false}); err != nil {
		t.Fatal(err)
	}
	if updateBody["certResolver"] != "custom" || len(updateBody) != 1 {
		t.Fatalf("update sent unchanged fields: %#v", updateBody)
	}
	if _, err := Call(context.Background(), c, d.Read, v, nil); err != nil {
		t.Fatal(err)
	}
	path, err := Path(d.Delete.Path, v)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Do(context.Background(), d.Delete.Method, path, nil, nil); err != nil {
		t.Fatal(err)
	}
}
