package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golangTestTask/internal/service"
	service_mocks "golangTestTask/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetBalance(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockWallet, address string, balance float64, err error)

	tests := []struct {
		name                 string
		address              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Success",
			address: "addr1",
			mockBehavior: func(s *service_mocks.MockWallet, address string, balance float64, err error) {
				s.EXPECT().GetWalletBalance(address).Return(balance, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"address":"addr1","balance":100.5}` + "\n",
		},
		{
			name:    "Wallet Not Found",
			address: "unknown",
			mockBehavior: func(s *service_mocks.MockWallet, address string, balance float64, err error) {
				s.EXPECT().GetWalletBalance(address).Return(0.0, errors.New("wallet not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "wallet not found\n",
		},
		{
			name:                 "Empty Address",
			address:              "",
			mockBehavior:         nil,
			expectedStatusCode:   http.StatusMovedPermanently,
			expectedResponseBody: "<a href=\"/api/wallet/balance\">Moved Permanently</a>.\n\n",
		},
		{
			name:                 "Long Address",
			address:              "this_is_a_very_long_wallet_address_that_exceeds_the_maximum_allowed_length_of_64_characters",
			mockBehavior:         nil,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "too long address\n",
		},
		{
			name:    "Service Error",
			address: "addr1",
			mockBehavior: func(s *service_mocks.MockWallet, address string, balance float64, err error) {
				s.EXPECT().GetWalletBalance(address).Return(0.0, errors.New("database error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			walletMock := service_mocks.NewMockWallet(c)
			if tt.mockBehavior != nil {
				tt.mockBehavior(walletMock, tt.address, 100.5, nil)
			}

			services := &service.Service{Wallet: walletMock}
			handler := NewHandler(services)

			r := http.NewServeMux()
			r.HandleFunc("GET /api/wallet/{address}/balance", handler.GetBalance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/wallet/"+tt.address+"/balance", nil)

			if req.URL.Path != "/api/wallet//balance" {
				req = req.WithContext(withPathValue(req.Context(), "address", tt.address))
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func withPathValue(parent context.Context, key, value string) context.Context {
	return context.WithValue(parent, contextKey(key), value)
}

type contextKey string
