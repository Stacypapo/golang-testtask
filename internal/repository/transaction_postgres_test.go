package repository

import (
	"errors"
	"testing"

	"golangTestTask/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   models.Transaction
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("from1", "to1", 10.5).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: models.Transaction{
				From:   "from1",
				To:     "to1",
				Amount: 10.5,
			},
		},
		{
			name: "Empty Fields",
			mock: func() {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("", "to1", 10.5).
					WillReturnError(errors.New("empty from address"))
			},
			input: models.Transaction{
				From:   "",
				To:     "to1",
				Amount: 10.5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.Create(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionPostgres_Getlast(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []models.Transaction
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "from_address", "to_address", "amount"}).
					AddRow(1, "from1", "to1", 10.5).
					AddRow(2, "from2", "to2", 20.0)

				mock.ExpectQuery("SELECT id, from_address, to_address, amount FROM transactions ORDER BY id DESC LIMIT \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			input: 2,
			want: []models.Transaction{
				{ID: 1, From: "from1", To: "to1", Amount: 10.5},
				{ID: 2, From: "from2", To: "to2", Amount: 20.0},
			},
		},
		{
			name: "Empty Result",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "from_address", "to_address", "amount"})

				mock.ExpectQuery("SELECT id, from_address, to_address, amount FROM transactions ORDER BY id DESC LIMIT \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			input: 2,
			want:  []models.Transaction{},
		},
		{
			name: "Database Error",
			mock: func() {
				mock.ExpectQuery("SELECT id, from_address, to_address, amount FROM transactions ORDER BY id DESC LIMIT \\$1").
					WithArgs(2).
					WillReturnError(errors.New("db error"))
			},
			input:   2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.Getlast(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err == nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
