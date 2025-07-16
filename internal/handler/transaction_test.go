package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golangTestTask/internal/models"
	"golangTestTask/internal/service"
	service_mocks "golangTestTask/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Send(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockTransaction, req models.Transaction)

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         models.Transaction
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Success",
			inputBody: `{"from": "addr1", "to": "addr2", "amount": 10.5}`,
			inputRequest: models.Transaction{
				From:   "addr1",
				To:     "addr2",
				Amount: 10.5,
			},
			mockBehavior: func(s *service_mocks.MockTransaction, req models.Transaction) {
				s.EXPECT().TransferFunds(req.From, req.To, req.Amount).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"success"}` + "\n",
		},
		{
			name:                 "Invalid JSON",
			inputBody:            `{"from": "addr1", "to": "addr2", "amount": "invalid"}`,
			inputRequest:         models.Transaction{},
			mockBehavior:         func(s *service_mocks.MockTransaction, req models.Transaction) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Invalid request body\n",
		},
		{
			name:                 "Missing Fields",
			inputBody:            `{"from": "", "to": "addr2", "amount": 10.5}`,
			inputRequest:         models.Transaction{},
			mockBehavior:         func(s *service_mocks.MockTransaction, req models.Transaction) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Missing required fields or invalid amount\n",
		},
		{
			name:      "Insufficient Funds",
			inputBody: `{"from": "addr1", "to": "addr2", "amount": 10.5}`,
			inputRequest: models.Transaction{
				From:   "addr1",
				To:     "addr2",
				Amount: 10.5,
			},
			mockBehavior: func(s *service_mocks.MockTransaction, req models.Transaction) {
				s.EXPECT().TransferFunds(req.From, req.To, req.Amount).Return(errors.New("insufficient funds"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "insufficient funds\n",
		},
		{
			name:      "Wallet Not Found",
			inputBody: `{"from": "addr1", "to": "addr2", "amount": 10.5}`,
			inputRequest: models.Transaction{
				From:   "addr1",
				To:     "addr2",
				Amount: 10.5,
			},
			mockBehavior: func(s *service_mocks.MockTransaction, req models.Transaction) {
				s.EXPECT().TransferFunds(req.From, req.To, req.Amount).Return(errors.New("sender wallet not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "sender wallet not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			transactionMock := service_mocks.NewMockTransaction(c)
			tt.mockBehavior(transactionMock, tt.inputRequest)

			services := &service.Service{Transaction: transactionMock}
			handler := NewHandler(services)

			r := http.NewServeMux()
			r.HandleFunc("/api/send", handler.Send)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/send", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetLastTransactions(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockTransaction, count int)

	tests := []struct {
		name                 string
		queryParam           string
		inputCount           int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Success",
			queryParam: "5",
			inputCount: 5,
			mockBehavior: func(s *service_mocks.MockTransaction, count int) {
				s.EXPECT().GetLastTransactions(count).Return([]models.Transaction{
					{ID: 1, From: "addr1", To: "addr2", Amount: 10.5},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"from":"addr1","to":"addr2","amount":10.5}]` + "\n",
		},
		{
			name:                 "Missing Count",
			queryParam:           "",
			inputCount:           0,
			mockBehavior:         func(s *service_mocks.MockTransaction, count int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Count parameter is required\n",
		},
		{
			name:                 "Invalid Count",
			queryParam:           "invalid",
			inputCount:           0,
			mockBehavior:         func(s *service_mocks.MockTransaction, count int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Count must be a positive integer\n",
		},
		{
			name:       "Empty Result",
			queryParam: "5",
			inputCount: 5,
			mockBehavior: func(s *service_mocks.MockTransaction, count int) {
				s.EXPECT().GetLastTransactions(count).Return([]models.Transaction{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]` + "\n",
		},
		{
			name:       "Service Error",
			queryParam: "5",
			inputCount: 5,
			mockBehavior: func(s *service_mocks.MockTransaction, count int) {
				s.EXPECT().GetLastTransactions(count).Return(nil, errors.New("database error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			transactionMock := service_mocks.NewMockTransaction(c)
			tt.mockBehavior(transactionMock, tt.inputCount)

			services := &service.Service{Transaction: transactionMock}
			handler := NewHandler(services)

			r := http.NewServeMux()
			r.HandleFunc("/api/transactions", handler.GetLast)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/transactions?count="+tt.queryParam, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
