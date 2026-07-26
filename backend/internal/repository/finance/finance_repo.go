package finance

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FinanceRepository interface {
	CreateFund(fund *domain.Fund) error
	GetFunds() ([]domain.Fund, error)
	GetFundByID(id uint) (*domain.Fund, error)
	CreateTransaction(tx *domain.Transaction) error
	GetTransactions(page int, limit int) ([]domain.Transaction, int64, error)
	GetTransactionsByFundID(fundID uint) ([]domain.Transaction, error)
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
	if err := r.db.Find(&funds).Error; err != nil {
		return nil, err
	}
	return funds, nil
}

func (r *financeRepository) GetFundByID(id uint) (*domain.Fund, error) {
	var fund domain.Fund
	if err := r.db.First(&fund, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fund, nil
}

func (r *financeRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		// 1. Get the fund and lock it for update
		var fund domain.Fund
		if err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fund, tx.FundID).Error; err != nil {
			return err
		}

		// 2. Validate sufficient balance if DEBIT
		if tx.Type == "DEBIT" {
			if fund.CurrentBalance < tx.Amount {
				return errors.New("insufficient fund balance")
			}
			fund.CurrentBalance -= tx.Amount
		} else if tx.Type == "CREDIT" {
			fund.CurrentBalance += tx.Amount
		} else {
			return errors.New("invalid transaction type")
		}

		// 3. Update the fund balance
		if err := dbTx.Model(&fund).Update("current_balance", fund.CurrentBalance).Error; err != nil {
			return err
		}

		// 4. Create the transaction record
		if err := dbTx.Create(tx).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *financeRepository) GetTransactions(page int, limit int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}
	
	return transactions, total, nil
}

func (r *financeRepository) GetTransactionsByFundID(fundID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.Where("fund_id = ?", fundID).Order("created_at desc").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
