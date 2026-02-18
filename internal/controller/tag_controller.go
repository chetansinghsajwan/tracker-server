package controller

import (
	"encoding/json"
	"net/http"

	"tracker-server/internal/service"
)

type TagController struct {
	tagService *service.TagService
}

func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

type createTagRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (c *TagController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := c.tagService.Create(r.Context(), req.UserID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

func (c *TagController) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	tags, err := c.tagService.ListByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tags)
}
