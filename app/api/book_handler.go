package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/theotruvelot/books/books"
	"github.com/theotruvelot/books/model"
)

type BookHandler struct {
	repository *books.BookRepository
	validate   *validator.Validate
}

func NewBookHandler(repo *books.BookRepository) *BookHandler {
	return &BookHandler{
		repository: repo,
		validate:   validator.New(),
	}
}

func (h *BookHandler) RegisterRoutes(r chi.Router) {
	r.Route("/books", func(r chi.Router) {
		r.Get("/", h.ListBooks)
		r.Post("/", h.CreateBook)
		r.Get("/{id}", h.GetBook)
		r.Put("/{id}", h.UpdateBook)
		r.Delete("/{id}", h.DeleteBook)
	})
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.repository.CreateBook(ctx, &book); err != nil {
		fmt.Println("Error creating book:", err)
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	book, err := h.repository.GetBook(ctx, id)
	if err != nil {
		http.Error(w, "Failed to retrieve book", http.StatusInternalServerError)
		return
	}

	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	book.ID = id
	if err := h.validate.Struct(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.repository.UpdateBook(ctx, &book); err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.repository.DeleteBook(ctx, id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	ctx := r.Context()
	books, err := h.repository.ListBooks(ctx, page, pageSize)
	if err != nil {
		fmt.Println("Error listing books:", err)
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
