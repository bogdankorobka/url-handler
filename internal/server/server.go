package server

import (
	"context"
	"log"
	"net/http"

	"github.com/bogdankorobka/url-handler/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title        URL Handler API
// @BasePath     /api/v1

type APIServer struct {
	server *http.Server
	host   string
}

func NewServer(host string) *APIServer {
	return &APIServer{host: host}
}

func (s *APIServer) Start() {
	s.server = &http.Server{
		Addr:    s.host,
		Handler: s.configureRouter(),
	}

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *APIServer) Stop(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) configureRouter() http.Handler {
	r := chi.NewRouter()

	// swagger routes
	docs.SwaggerInfo.Host = s.host
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json")),
	)

	r.Route("/api/v1", func(r chi.Router) {
		// common middlewares
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		// requests limiter
		r.Use(Limit(3))

		r.Post("/url-handler", UrlHandler())

	})

	return r
}
