package handler

import (
	"golangTestTask/internal/service"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *service.Service
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes инициализирует маршруты HTTP для обработчика Handler и возвращает настроенный мультиплексор (*http.ServeMux).
func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/send", h.Send)
	router.HandleFunc("GET /api/transactions", h.GetLast)
	router.HandleFunc("GET /api/wallet/{address}/balance", h.GetBalance)
	router.Handle("/swagger/", httpSwagger.WrapHandler)
	return router
}
