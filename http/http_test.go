package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = new(http.Client)

func Test指定したパスにアクセスしたときにJSONを返すことを確認するテスト(t *testing.T) {
	type args struct {
		req  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "/api/posts/にアクセスしたときに指定したJSONを返すことを確認すること",
			args: args{
				req:  "GET",
				path: "/api/posts",
			},
			want: "{\"id\":0,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\",\"title\":\"Test desu\",\"body\":\"Hello World\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			s.Router()
			testServer := httptest.NewServer(s.service)
			defer testServer.Close()
			path := testServer.URL + tt.args.path
			req, _ := http.NewRequest(tt.args.req, path, nil)

			resp, _ := client.Do(req)
			respBody, _ := ioutil.ReadAll(resp.Body)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tt.want, string(respBody))
		})
	}
}
