package handler

import (
	"encoding/json"
	"golangTestTask/internal/models"
	"net/http"
	"strconv"
)

// Send
// @Summary Отправить денежные средства
// @Description Переводит денежные средства с одного кошелька на другой
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Данные транзакции"
// @Success 200 {object} models.StatusResponse "Status"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 404 {string} string "Wallet not found"
// @Router /api/send [post]
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.From == "" || req.To == "" || req.Amount <= 0 {
		http.Error(w, "Missing required fields or invalid amount", http.StatusBadRequest)
		return
	}

	if err := h.services.TransferFunds(req.From, req.To, req.Amount); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "insufficient funds" {
			status = http.StatusBadRequest
		} else if err.Error() == "sender wallet not found" || err.Error() == "recipient wallet not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.StatusResponse{
		Status:  "success",
		Message: "Transaction completed",
	})
}

// GetLastTransactions возвращает N последних транзакций
// @Summary Получить последние транзакции
// @Description Возвращает N последних по времени переводов средств
// @Produce json
// @Param count query int false "Количество транзакций"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Server error"
// @Router /api/transactions [get]
func (h *Handler) GetLast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		http.Error(w, "Count parameter is required", http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		http.Error(w, "Count must be a positive integer", http.StatusBadRequest)
		return
	}

	transactions, err := h.services.GetLastTransactions(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
