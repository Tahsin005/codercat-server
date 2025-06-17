package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tahsin005/codercat-server/domain"
	"github.com/tahsin005/codercat-server/service"
)

type BlogHandler struct {
	service service.BlogService
}

func NewBlogHandler(service service.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog domain.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateBlog(r.Context(), &blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	blog, err := h.service.GetBlogByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var blog domain.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateBlog(r.Context(), id, &blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.DeleteBlog(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.service.GetAllBlogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetFeaturedBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.service.GetFeaturedBlogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetRecentBlogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 3
	}
	blogs, err := h.service.GetRecentBlogs(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetBlogsByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	blogs, err := h.service.GetBlogsByCategory(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) SearchBlogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	blogs, err := h.service.SearchBlogs(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetRelatedBlogs(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 2
	}
	blogs, err := h.service.GetRelatedBlogs(r.Context(), id, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *BlogHandler) GetPopularCategories(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 5 // default to top 5 categories
	}

	categories, err := h.service.GetPopularCategories(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *BlogHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/blogs/featured", h.GetFeaturedBlogs).Methods("GET")
	router.HandleFunc("/blogs/recent", h.GetRecentBlogs).Methods("GET")
	router.HandleFunc("/blogs/search", h.SearchBlogs).Methods("GET")
	router.HandleFunc("/categories", h.GetCategories).Methods("GET")
	router.HandleFunc("/categories/popular", h.GetPopularCategories).Methods("GET")

	router.HandleFunc("/blogs/related/{id}", h.GetRelatedBlogs).Methods("GET")
	router.HandleFunc("/blogs/{id}", h.GetBlog).Methods("GET")
	router.HandleFunc("/blogs/{id}", h.UpdateBlog).Methods("PUT")
	router.HandleFunc("/blogs/{id}", h.DeleteBlog).Methods("DELETE")

	router.HandleFunc("/blogs", h.CreateBlog).Methods("POST")
	router.HandleFunc("/blogs", h.GetAllBlogs).Methods("GET")
	router.HandleFunc("/blogs/category/{category}", h.GetBlogsByCategory).Methods("GET")
}
