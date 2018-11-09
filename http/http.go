package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/keiya01/eserver/controller"
)

// Server Server
type Server struct {
	service *chi.Mux
}

// New Server構造体のコンストラクタ
func NewServer() *Server {
	return &Server{
		service: chi.NewRouter(),
	}
}

// Router ルーティング設定
func (s *Server) Router() {
	p := controller.NewPostController()

	s.service.Route("/api", func(api chi.Router) {
		// api.Use(Auth("db connection"))
		api.Route("/posts", func(posts chi.Router) {
			posts.Get("/", p.Index)
		})
	})
}

func (s *Server) Start(port string) {
	err := http.ListenAndServe(port, s.service)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
