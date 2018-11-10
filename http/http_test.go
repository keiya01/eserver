package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
