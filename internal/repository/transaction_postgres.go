package repository

import (
	"database/sql"
	"fmt"
	"golangTestTask/internal/models"
)

type TransactionPostgres struct {
	db *sql.DB
}

// NewTransactionPostgres создает новый экземпляр TransactionPostgres.
func NewTransactionPostgres(db *sql.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

// Create сохраняет новую транзакцию в БД PostgreSQL.
func (r *TransactionPostgres) Create(transaction models.Transaction) error {
	query := `INSERT INTO transactions (from_address, to_address, amount) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

// Getlast возвращает count последних транзакций из БД PostgreSQL, отсортированных по ID в порядке убывания.
func (r *TransactionPostgres) Getlast(count int) ([]models.Transaction, error) {
	query := `SELECT id, from_address, to_address, amount FROM transactions ORDER BY id DESC LIMIT $1`
	transactions := make([]models.Transaction, 0)

	rows, err := r.db.Query(query, count)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.From, &t.To, &t.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return transactions, nil
}
