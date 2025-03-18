package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/theotruvelot/books/books"
)

type Server struct {
	router      chi.Router
	bookRepo    *books.BookRepository
	bookHandler *BookHandler
}

func NewServer(bookRepo *books.BookRepository) *Server {
	s := &Server{
		router:   chi.NewRouter(),
		bookRepo: bookRepo,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
}

func (s *Server) setupRoutes() {
	s.bookHandler = NewBookHandler(s.bookRepo)
	s.bookHandler.RegisterRoutes(s.router)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
