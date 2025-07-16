package repository

import (
	"database/sql"
	"errors"
	"testing"

	"golangTestTask/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestWalletPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewWalletPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   *models.Wallet
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("INSERT INTO wallets").
					WithArgs("addr1", 100.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &models.Wallet{
				Address: "addr1",
				Balance: 100.0,
			},
			wantErr: false,
		},
		{
			name: "Duplicate Address",
			mock: func() {
				mock.ExpectExec("INSERT INTO wallets").
					WithArgs("addr1", 100.0).
					WillReturnError(errors.New("duplicate key"))
			},
			input: &models.Wallet{
				Address: "addr1",
				Balance: 100.0,
			},
			wantErr: true,
		},
		{
			name: "Empty Address",
			mock: func() {
				mock.ExpectExec("INSERT INTO wallets").
					WithArgs("", 100.0).
					WillReturnError(errors.New("empty address"))
			},
			input: &models.Wallet{
				Address: "",
				Balance: 100.0,
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

func TestWalletPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewWalletPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   *models.Wallet
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("UPDATE wallets").
					WithArgs(150.0, "addr1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: &models.Wallet{
				Address: "addr1",
				Balance: 150.0,
			},
			wantErr: false,
		},
		{
			name: "Wallet Not Found",
			mock: func() {
				mock.ExpectExec("UPDATE wallets").
					WithArgs(150.0, "unknown").
					WillReturnError(errors.New("wallet not found"))
			},
			input: &models.Wallet{
				Address: "unknown",
				Balance: 150.0,
			},
			wantErr: true,
		},
		{
			name: "Database Error",
			mock: func() {
				mock.ExpectExec("UPDATE wallets").
					WithArgs(150.0, "addr1").
					WillReturnError(errors.New("db error"))
			},
			input: &models.Wallet{
				Address: "addr1",
				Balance: 150.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.Update(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletPostgres_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewWalletPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *models.Wallet
		wantErr error
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"address", "balance"}).
					AddRow("addr1", 100.0)
				mock.ExpectQuery("SELECT address, balance FROM wallets").
					WithArgs("addr1").
					WillReturnRows(rows)
			},
			input: "addr1",
			want: &models.Wallet{
				Address: "addr1",
				Balance: 100.0,
			},
			wantErr: nil,
		},
		{
			name: "Wallet Not Found",
			mock: func() {
				mock.ExpectQuery("SELECT address, balance FROM wallets").
					WithArgs("unknown").
					WillReturnError(sql.ErrNoRows)
			},
			input:   "unknown",
			want:    nil,
			wantErr: ErrWalletNotFound,
		},
		{
			name: "Database Error",
			mock: func() {
				mock.ExpectQuery("SELECT address, balance FROM wallets").
					WithArgs("addr1").
					WillReturnError(errors.New("db error"))
			},
			input:   "addr1",
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.Get(tt.input)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletPostgres_Existence(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewWalletPostgres(db)

	tests := []struct {
		name     string
		mock     func()
		expected bool
	}{
		{
			name: "Exists",
			mock: func() {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)
				mock.ExpectQuery("SELECT EXISTS").
					WillReturnRows(rows)
			},
			expected: true,
		},
		{
			name: "Not Exists",
			mock: func() {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(false)
				mock.ExpectQuery("SELECT EXISTS").
					WillReturnRows(rows)
			},
			expected: false,
		},
		{
			name: "Database Error",
			mock: func() {
				mock.ExpectQuery("SELECT EXISTS").
					WillReturnError(errors.New("db error"))
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result := repo.Existence()
			assert.Equal(t, tt.expected, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
