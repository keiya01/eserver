package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var client = new(http.Client)

func Test指定したパスにアクセスしたときにJSONを返すことを確認するテスト(t *testing.T) {
	s := NewServer()
	s.Router()
	testServer := httptest.NewServer(s.service)
	defer testServer.Close()
	path := testServer.URL + "/api/posts"
	req, _ := http.NewRequest("GET", path, nil)

	resp, _ := client.Do(req)

	if  resp.StatusCode != 200 {
		t.Errorf("HTTP Test : get = %d", resp.StatusCode)
	}
}
