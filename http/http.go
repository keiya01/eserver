package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/keiya01/eserver/auth"
	"github.com/keiya01/eserver/controller"
)

// Server Server
type Server struct {
	*chi.Mux
}

// NewServer Server構造体のコンストラクタ
func NewServer() *Server {
	return &Server{
		Mux: chi.NewRouter(),
	}
}

// Router ルーティング設定
func (s *Server) Router() {
	cors := corsNew()
	s.Use(cors.Handler)
	s.Use(middleware.RequestID)
	s.Use(middleware.Logger)
	s.Use(middleware.URLFormat)
	s.Use(render.SetContentType(render.ContentTypeJSON))
	s.Use(auth.JWTAuthentication) // ここでJWTの認証を行う

	p := controller.PostController{}
	u := controller.UserController{}

	s.Route("/api", func(api chi.Router) {
		api.Route("/posts", func(posts chi.Router) {
			posts.Get("/{id}", p.Show)
			posts.Post("/create", p.Create)
			posts.Put("/{id}/update", p.Update)
			posts.Delete("/{id}/delete", p.Delete)
		})
		api.Route("/users", func(users chi.Router) {
			users.Post("/login", u.Login)
			users.Post("/create", u.Create)
			users.Get("/{id}", u.Show)
			users.Put("/{id}/update", u.Update)
			users.Delete("/{id}/delete", u.Delete)
		})
	})
}

// Start 指定したPortでローカルサーバーを起動する
func (s *Server) Start(port string) {
	err := http.ListenAndServe(port, s)
	if err != nil {
		log.Println("ListenAndServe:", err)
	}
}

func corsNew() *cors.Cors {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return cors
}
