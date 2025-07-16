package service

import (
	"golangTestTask/internal/models"
	"golangTestTask/internal/repository"
)

type Wallet interface {
	// CreateWallet создает новый кошелек.
	CreateWallet(models.Wallet) error
	// GetWalletBalance возвращает баланс кошелька по его адресу
	GetWalletBalance(address string) (float64, error)
	// CreateRandomWallets создает count кошельков со случайными адресами и balance у.е. на них.
	CreateRandomWallets(count int, balance float64) error
	// BaseWallets создает count кошельков со случайными адресами и balance у.е. на них если они еще не созданы.
	BaseWallets(count int, balance float64) error
}

type Transaction interface {
	// TransferFunds переводит средства между кошельками
	TransferFunds(from string, to string, amount float64) error
	// GetLastTransactions возвращает последние count транзакций.
	GetLastTransactions(count int) ([]models.Transaction, error)
}

type Service struct {
	Wallet
	Transaction
}

// NewService создает новый экземпляр Service.
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Wallet:      NewWalletService(repo.Wallet),
		Transaction: NewTransactionService(repo.Transaction, repo.Wallet),
	}
}
