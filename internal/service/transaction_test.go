package service

import (
	"errors"
	"testing"

	"golangTestTask/internal/models"
	"golangTestTask/internal/repository"
	repository_mocks "golangTestTask/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransactionService_TransferFunds(t *testing.T) {
	type mockBehavior struct {
		getFrom    func(r *repository_mocks.MockWallet, from string, balance float64)
		getTo      func(r *repository_mocks.MockWallet, to string, balance float64)
		updateFrom func(r *repository_mocks.MockWallet, wallet *models.Wallet)
		updateTo   func(r *repository_mocks.MockWallet, wallet *models.Wallet)
		createTx   func(r *repository_mocks.MockTransaction, tx models.Transaction)
	}

	tests := []struct {
		name         string
		from         string
		to           string
		amount       float64
		mockBehavior mockBehavior
		wantErr      bool
		expectedErr  string
	}{
		{
			name:   "successful transfer",
			from:   "addr1",
			to:     "addr2",
			amount: 10.5,
			mockBehavior: mockBehavior{
				getFrom: func(r *repository_mocks.MockWallet, from string, balance float64) {
					r.EXPECT().Get(from).Return(&models.Wallet{
						Address: from,
						Balance: 100.0,
					}, nil)
				},
				getTo: func(r *repository_mocks.MockWallet, to string, balance float64) {
					r.EXPECT().Get(to).Return(&models.Wallet{
						Address: to,
						Balance: 50.0,
					}, nil)
				},
				updateFrom: func(r *repository_mocks.MockWallet, wallet *models.Wallet) {
					r.EXPECT().Update(wallet).Return(nil)
				},
				updateTo: func(r *repository_mocks.MockWallet, wallet *models.Wallet) {
					r.EXPECT().Update(wallet).Return(nil)
				},
				createTx: func(r *repository_mocks.MockTransaction, tx models.Transaction) {
					r.EXPECT().Create(tx).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name:   "sender not found",
			from:   "unknown",
			to:     "addr2",
			amount: 10.5,
			mockBehavior: mockBehavior{
				getFrom: func(r *repository_mocks.MockWallet, from string, balance float64) {
					r.EXPECT().Get(from).Return(nil, repository.ErrWalletNotFound)
				},
			},
			wantErr:     true,
			expectedErr: "sender wallet not found",
		},
		{
			name:   "recipient not found",
			from:   "addr1",
			to:     "unknown",
			amount: 10.5,
			mockBehavior: mockBehavior{
				getFrom: func(r *repository_mocks.MockWallet, from string, balance float64) {
					r.EXPECT().Get(from).Return(&models.Wallet{
						Address: from,
						Balance: 100.0,
					}, nil)
				},
				getTo: func(r *repository_mocks.MockWallet, to string, balance float64) {
					r.EXPECT().Get(to).Return(nil, repository.ErrWalletNotFound)
				},
			},
			wantErr:     true,
			expectedErr: "recipient wallet not found",
		},
		{
			name:   "insufficient funds",
			from:   "addr1",
			to:     "addr2",
			amount: 150.0,
			mockBehavior: mockBehavior{
				getFrom: func(r *repository_mocks.MockWallet, from string, balance float64) {
					r.EXPECT().Get(from).Return(&models.Wallet{
						Address: from,
						Balance: 100.0,
					}, nil)
				},
			},
			wantErr:     true,
			expectedErr: "insufficient funds",
		},
		{
			name:   "update sender failed",
			from:   "addr1",
			to:     "addr2",
			amount: 10.5,
			mockBehavior: mockBehavior{
				getFrom: func(r *repository_mocks.MockWallet, from string, balance float64) {
					r.EXPECT().Get(from).Return(&models.Wallet{
						Address: from,
						Balance: 100.0,
					}, nil)
				},
				getTo: func(r *repository_mocks.MockWallet, to string, balance float64) {
					r.EXPECT().Get(to).Return(&models.Wallet{
						Address: to,
						Balance: 50.0,
					}, nil)
				},
				updateFrom: func(r *repository_mocks.MockWallet, wallet *models.Wallet) {
					r.EXPECT().Update(wallet).Return(errors.New("update failed"))
				},
			},
			wantErr:     true,
			expectedErr: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			walletRepo := repository_mocks.NewMockWallet(ctrl)
			txRepo := repository_mocks.NewMockTransaction(ctrl)

			if tt.mockBehavior.getFrom != nil {
				tt.mockBehavior.getFrom(walletRepo, tt.from, 100.0)
			}
			if tt.mockBehavior.getTo != nil {
				tt.mockBehavior.getTo(walletRepo, tt.to, 50.0)
			}
			if tt.mockBehavior.updateFrom != nil {
				tt.mockBehavior.updateFrom(walletRepo, &models.Wallet{
					Address: tt.from,
					Balance: 100.0 - tt.amount,
				})
			}
			if tt.mockBehavior.updateTo != nil {
				tt.mockBehavior.updateTo(walletRepo, &models.Wallet{
					Address: tt.to,
					Balance: 50.0 + tt.amount,
				})
			}
			if tt.mockBehavior.createTx != nil {
				tt.mockBehavior.createTx(txRepo, models.Transaction{
					From:   tt.from,
					To:     tt.to,
					Amount: tt.amount,
				})
			}

			service := NewTransactionService(txRepo, walletRepo)
			err := service.TransferFunds(tt.from, tt.to, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransactionService_GetLastTransactions(t *testing.T) {
	tests := []struct {
		name           string
		count          int
		mockBehavior   func(r *repository_mocks.MockTransaction, count int)
		expectedResult []models.Transaction
		wantErr        bool
	}{
		{
			name:  "successful get last transactions",
			count: 2,
			mockBehavior: func(r *repository_mocks.MockTransaction, count int) {
				r.EXPECT().Getlast(count).Return([]models.Transaction{
					{From: "addr1", To: "addr2", Amount: 10.5},
					{From: "addr2", To: "addr1", Amount: 5.0},
				}, nil)
			},
			expectedResult: []models.Transaction{
				{From: "addr1", To: "addr2", Amount: 10.5},
				{From: "addr2", To: "addr1", Amount: 5.0},
			},
			wantErr: false,
		},
		{
			name:  "repository error",
			count: 2,
			mockBehavior: func(r *repository_mocks.MockTransaction, count int) {
				r.EXPECT().Getlast(count).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:  "empty result",
			count: 2,
			mockBehavior: func(r *repository_mocks.MockTransaction, count int) {
				r.EXPECT().Getlast(count).Return([]models.Transaction{}, nil)
			},
			expectedResult: []models.Transaction{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			txRepo := repository_mocks.NewMockTransaction(ctrl)
			tt.mockBehavior(txRepo, tt.count)

			service := NewTransactionService(txRepo, nil)
			result, err := service.GetLastTransactions(tt.count)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
