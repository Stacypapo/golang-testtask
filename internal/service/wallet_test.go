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

func TestWalletService_CreateWallet(t *testing.T) {
	tests := []struct {
		name        string
		wallet      models.Wallet
		mock        func(*repository_mocks.MockWallet, *models.Wallet)
		expectedErr error
	}{
		{
			name: "success",
			wallet: models.Wallet{
				Address: "addr1",
				Balance: 100.0,
			},
			mock: func(m *repository_mocks.MockWallet, w *models.Wallet) {
				m.EXPECT().Create(w).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			wallet: models.Wallet{
				Address: "addr1",
				Balance: 100.0,
			},
			mock: func(m *repository_mocks.MockWallet, w *models.Wallet) {
				m.EXPECT().Create(w).Return(errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository_mocks.NewMockWallet(ctrl)
			tt.mock(mockRepo, &tt.wallet)

			service := NewWalletService(mockRepo)
			err := service.CreateWallet(tt.wallet)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWalletService_GetWalletBalance(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		mock        func(*repository_mocks.MockWallet, string)
		expectedBal float64
		expectedErr error
	}{
		{
			name:    "success",
			address: "addr1",
			mock: func(m *repository_mocks.MockWallet, addr string) {
				m.EXPECT().Get(addr).Return(&models.Wallet{
					Address: addr,
					Balance: 100.0,
				}, nil)
			},
			expectedBal: 100.0,
			expectedErr: nil,
		},
		{
			name:    "wallet not found",
			address: "unknown",
			mock: func(m *repository_mocks.MockWallet, addr string) {
				m.EXPECT().Get(addr).Return(nil, repository.ErrWalletNotFound)
			},
			expectedBal: 0,
			expectedErr: repository.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository_mocks.NewMockWallet(ctrl)
			tt.mock(mockRepo, tt.address)

			service := NewWalletService(mockRepo)
			balance, err := service.GetWalletBalance(tt.address)

			assert.Equal(t, tt.expectedBal, balance)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestWalletService_CreateRandomWallets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository_mocks.NewMockWallet(ctrl)

	// Ожидаем 3 вызова Create
	mockRepo.EXPECT().Create(gomock.Any()).Times(3).Return(nil)

	service := NewWalletService(mockRepo)
	err := service.CreateRandomWallets(3, 100.0)

	assert.NoError(t, err)
}

func TestWalletService_BaseWallets(t *testing.T) {
	tests := []struct {
		name        string
		count       int
		balance     float64
		mock        func(*repository_mocks.MockWallet)
		expectedErr error
	}{
		{
			name:    "create new wallets",
			count:   3,
			balance: 100.0,
			mock: func(m *repository_mocks.MockWallet) {
				m.EXPECT().Existence().Return(false)
				m.EXPECT().Create(gomock.Any()).Times(3).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "wallets already exist",
			count:   3,
			balance: 100.0,
			mock: func(m *repository_mocks.MockWallet) {
				m.EXPECT().Existence().Return(true)
			},
			expectedErr: errors.New("wallets already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository_mocks.NewMockWallet(ctrl)
			tt.mock(mockRepo)

			service := NewWalletService(mockRepo)
			err := service.BaseWallets(tt.count, tt.balance)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
