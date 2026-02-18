package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tracker-server/internal/repo"
	"tracker-server/internal/service"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

type createAccountRequest struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

func (c *AccountController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := c.accountService.Create(r.Context(), req.UserID, req.Name, req.Type, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (c *AccountController) ListAndGet(w http.ResponseWriter, r *http.Request) {
	// Simple routing within handler: /accounts?user_id=... (List) OR /accounts/{id} (Get) via separate handler?
	// Given main.go routing limitations, I'll likely mount this on /accounts/

	idStr := ""
	if len(r.URL.Path) > len("/accounts/") {
		idStr = r.URL.Path[len("/accounts/"):]
	}

	if r.Method == http.MethodGet {
		if idStr != "" {
			c.Get(w, r, idStr)
		} else {
			c.List(w, r)
		}
		return
	}

	if r.Method == http.MethodPut && idStr != "" {
		c.Update(w, r, idStr)
		return
	}

	if r.Method == http.MethodDelete && idStr != "" {
		c.Delete(w, r, idStr)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (c *AccountController) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id query parameter", http.StatusBadRequest)
		return
	}

	accounts, err := c.accountService.ListByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (c *AccountController) Get(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	account, err := c.accountService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (c *AccountController) Update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var account repo.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	account.ID = id // Ensure ID matches URL

	if err := c.accountService.Update(r.Context(), &account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AccountController) Delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.accountService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
