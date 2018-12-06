package request

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
