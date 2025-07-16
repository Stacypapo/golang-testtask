package handler

import (
	"encoding/json"
	"golangTestTask/internal/models"
	"net/http"
)

// GetBalance возвращает баланс кошелька
// @Summary Получить баланс кошелька
// @Description Возвращает баланс по адресу кошелька
// @Produce json
// @Param address path string true "Адрес кошелька"
// @Success 200 {object} models.Wallet
// @Failure 400 {string} string "Invalid address"
// @Failure 404 {string} string "Wallet not found"
// @Failure 500 {string} string "Server error"
// @Router /api/wallet/{address}/balance [get]
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.PathValue("address")
	if len(address) >= 64 {
		http.Error(w, "too long address", http.StatusBadRequest)
		return
	}

	balance, err := h.services.GetWalletBalance(address)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "wallet not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	response := models.Wallet{
		Address: address,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
