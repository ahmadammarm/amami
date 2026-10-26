package finance

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionRow is a ledger row plus the id of the reversal entry that cancels it, if any.
type TransactionRow struct {
	domain.Transaction `gorm:"embedded"`
	ReversedByID       *uuid.UUID
}

// ReversalBuilder receives the locked original transaction and returns the
// entry that reverses it, or an error if it must not be reversed.
type ReversalBuilder func(original *domain.Transaction) (*domain.Transaction, error)

type FinanceRepository interface {
	CreateFund(fund *domain.Fund) error
	GetFunds() ([]domain.Fund, error)
	CreateTransaction(tx *domain.Transaction) error
	ReverseTransaction(originalID uuid.UUID, build ReversalBuilder) (*domain.Transaction, error)
	GetTransactions(page int, limit int) ([]TransactionRow, int64, error)
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) CreateFund(fund *domain.Fund) error {
	return r.db.Create(fund).Error
}

func (r *financeRepository) GetFunds() ([]domain.Fund, error) {
	var funds []domain.Fund
	if err := r.db.Order("id").Find(&funds).Error; err != nil {
		return nil, err
	}
	return funds, nil
}

func (r *financeRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		return PostTransaction(dbTx, tx)
	})
}

func (r *financeRepository) ReverseTransaction(originalID uuid.UUID, build ReversalBuilder) (*domain.Transaction, error) {
	var reversal *domain.Transaction
	// Lock order: the transaction row, then (inside PostTransaction) the fund row.
	err := r.db.Transaction(func(dbTx *gorm.DB) error {
		original, err := database.LockForUpdate[domain.Transaction](dbTx, originalID)
		if err != nil {
			return err
		}

		var existing int64
		if err := dbTx.Model(&domain.Transaction{}).
			Where("reversal_of_id = ?", originalID).Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return domain.ErrAlreadyReversed
		}

		built, err := build(original)
		if err != nil {
			return err
		}
		if err := PostTransaction(dbTx, built); err != nil {
			return err
		}
		reversal = built
		return nil
	})

	// The unique index on reversal_of_id catches a concurrent reversal that
	// slipped past the count check above.
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, domain.ErrAlreadyReversed
	}
	if err != nil {
		return nil, err
	}
	return reversal, nil
}

func (r *financeRepository) GetTransactions(page int, limit int) ([]TransactionRow, int64, error) {
	var total int64
	if err := r.db.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []TransactionRow
	offset := (page - 1) * limit
	err := r.db.Model(&domain.Transaction{}).
		Select("transactions.*, reversal.id AS reversed_by_id").
		Joins("LEFT JOIN transactions AS reversal ON reversal.reversal_of_id = transactions.id").
		Order("transactions.created_at desc").
		Offset(offset).Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
