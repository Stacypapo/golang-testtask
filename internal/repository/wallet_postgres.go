package repository

import (
	"database/sql"
	"errors"
	"golangTestTask/internal/models"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletPostgres struct {
	db *sql.DB
}

// NewWalletPostgres создает новый экземпляр WalletPostgres.
func NewWalletPostgres(db *sql.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

// Create сохраняет новый кошелек в БД PostgreSQL.
func (r *WalletPostgres) Create(wallet *models.Wallet) error {
	query := `INSERT INTO wallets (address, balance) VALUES ($1, $2)`
	_, err := r.db.Exec(query, wallet.Address, wallet.Balance)
	if err != nil {
		return err
	}
	return nil
}

// Update обновляет баланс кошелька по адресу в БД PostgreSQL.
func (r *WalletPostgres) Update(wallet *models.Wallet) error {
	query := `UPDATE wallets SET balance = $1 WHERE address = $2`
	_, err := r.db.Exec(query, wallet.Balance, wallet.Address)
	if err != nil {
		return err
	}
	return nil
}

// Get возвращает кошелек по адресу в БД PostgreSQL.
func (r *WalletPostgres) Get(address string) (*models.Wallet, error) {
	query := `SELECT address, balance FROM wallets WHERE address = $1`
	row := r.db.QueryRow(query, address)

	var wallet models.Wallet
	err := row.Scan(&wallet.Address, &wallet.Balance)
	if err == sql.ErrNoRows {
		return nil, ErrWalletNotFound
	}
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Existence проверяет существуют ли какие-либо кошельки в БД PostgreSQL.
func (r *WalletPostgres) Existence() bool {
	query := `SELECT EXISTS (SELECT 1 FROM wallets)`
	var exists bool
	err := r.db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
