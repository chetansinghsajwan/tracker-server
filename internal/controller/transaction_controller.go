package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"tracker-server/internal/repo"
	"tracker-server/internal/service"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

func (c *TransactionController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t repo.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.transactionService.Create(r.Context(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (c *TransactionController) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	filter := repo.TransactionFilter{}
	// Parse filters
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}
	if minAmount := r.URL.Query().Get("min_amount"); minAmount != "" {
		if f, err := strconv.ParseFloat(minAmount, 64); err == nil {
			filter.MinAmount = &f
		}
	}
	if maxAmount := r.URL.Query().Get("max_amount"); maxAmount != "" {
		if f, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			filter.MaxAmount = &f
		}
	}
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if id, err := strconv.ParseInt(categoryID, 10, 64); err == nil {
			filter.CategoryID = &id
		}
	}

	transactions, err := c.transactionService.ListByUserID(r.Context(), userID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
