package service

import (
	"errors"
	"golangTestTask/internal/models"
	"golangTestTask/internal/repository"
	"golangTestTask/pkg/utils"
)

type WalletService struct {
	repo repository.Wallet
}

// NewWalletService создает новый экземпляр WalletService.
func NewWalletService(repo repository.Wallet) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

// CreateWallet создает новый кошелек.
func (s *WalletService) CreateWallet(wallet models.Wallet) error {
	if err := s.repo.Create(&wallet); err != nil {
		return err
	}
	return nil
}

// GetWalletBalance возвращает баланс кошелька по его адресу
func (s *WalletService) GetWalletBalance(address string) (float64, error) {
	wallet, err := s.repo.Get(address)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// GetAllWallets возвращает все кошельки в базе данных
func (s *WalletService) GetAllWallets() ([]models.Wallet, error) {
	var wallets []models.Wallet
	wallets, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateRandomWallets создает count кошельков со случайными адресами и balance у.е. на них.
func (s *WalletService) CreateRandomWallets(count int, balance float64) error {
	for i := 0; i < count; i++ {
		s.CreateWallet(models.Wallet{
			Address: utils.GenerateAddress(),
			Balance: balance,
		})
	}
	return nil
}

// BaseWallets создает count кошельков со случайными адресами и balance у.е. на них если они еще не созданы.
func (s *WalletService) BaseWallets(count int, balance float64) error {
	if s.repo.Existence() {
		return errors.New("wallets already exists")
	}
	s.CreateRandomWallets(count, balance)
	return nil
}
