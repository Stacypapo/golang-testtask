package service

import (
	"errors"
	"fmt"
	"golangTestTask/internal/models"
	"golangTestTask/internal/repository"
)

type TransactionService struct {
	transaction_repo repository.Transaction
	wallet_repo      repository.Wallet
}

// NewTransactionService создает новый экземпляр TransactionService.
func NewTransactionService(transaction_repo repository.Transaction, wallet_repo repository.Wallet) *TransactionService {
	return &TransactionService{
		transaction_repo: transaction_repo,
		wallet_repo:      wallet_repo,
	}
}

// TransferFunds переводит amount средств из кошелька from на кошелек to.
func (s *TransactionService) TransferFunds(from string, to string, amount float64) error {
	var wallet_from, wallet_to *models.Wallet
	wallet_from, err := s.wallet_repo.Get(from)
	if err != nil {
		return fmt.Errorf("sender %w", err)
	}
	if wallet_from.Balance < amount {
		return errors.New("insufficient funds")
	}
	wallet_to, err = s.wallet_repo.Get(to)
	if err != nil {
		return fmt.Errorf("recipient %w", err)
	}

	wallet_from.Balance -= amount
	wallet_to.Balance += amount
	if err := s.wallet_repo.Update(wallet_from); err != nil {
		return err
	}
	if err := s.wallet_repo.Update(wallet_to); err != nil {
		return err
	}
	s.transaction_repo.Create(models.Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	})
	return nil
}

// GetLastTransactions возвращает последние count транзакций.
func (s *TransactionService) GetLastTransactions(count int) ([]models.Transaction, error) {
	var last_transactions []models.Transaction
	last_transactions, err := s.transaction_repo.Getlast(count)
	if err != nil {
		return nil, err
	}
	return last_transactions, nil
}
