package httpapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"task051-base62/internal/httpapi"
)

func TestBug2AllocBatchHTTPReturnsResults(t *testing.T) {
	ts := httptest.NewServer(httpapi.NewServer().Handler())
	defer ts.Close()

	body, _ := json.Marshal(map[string]any{"sources": []string{"A", "B"}})
	resp, err := http.Post(ts.URL+"/alloc-batch", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("alloc-batch request: %v", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body=%s, want 200", resp.StatusCode, data)
	}
	var got struct {
		OK      bool `json:"ok"`
		Results []struct {
			Source string `json:"source"`
			Code   string `json:"code"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v (body=%s)", err, data)
	}
	if !got.OK {
		t.Fatalf("ok=false, want true; body=%s", data)
	}
	if len(got.Results) != 2 {
		t.Fatalf("results = %+v (len %d), want 2 entries", got.Results, len(got.Results))
	}
	if got.Results[0].Code != "0" || got.Results[1].Code != "1" {
		t.Fatalf("codes = %+v, want [0 1]", got.Results)
	}
}
