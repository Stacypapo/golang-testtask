package repository

import (
	"database/sql"
	"golangTestTask/internal/models"
)

type Wallet interface {
	// Create сохраняет новый кошелек в БД.
	Create(wallet *models.Wallet) error
	// Update обновляет баланс кошелька по адресу.
	Update(wallet *models.Wallet) error
	// Get возвращает кошелек по адресу.
	Get(address string) (*models.Wallet, error)
	// Existence проверяет существуют ли какие-либо кошельки в БД.
	Existence() bool
}

type Transaction interface {
	// Create сохраняет новую транзакцию в БД.
	Create(transaction models.Transaction) error
	// Getlast возвращает count последних транзакций из БД.
	Getlast(count int) ([]models.Transaction, error)
}

type Repository struct {
	Wallet
	Transaction
}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Wallet:      NewWalletPostgres(db),
		Transaction: NewTransactionPostgres(db),
	}
}
