package finance

import (
	"errors"
	"fmt"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostTransaction appends tx to the ledger and applies it to the fund balance.
// It must run inside a database transaction: the fund row is locked so
// concurrent postings to the same fund are serialized. Every ledger write in
// the system (manual entries, reversals, zakat) goes through this function.
func PostTransaction(dbTx *gorm.DB, tx *domain.Transaction) error {
	var fund domain.Fund
	if err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fund, tx.FundID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrFundNotFound
		}
		return err
	}

	switch tx.Type {
	case domain.TransactionCredit:
		fund.CurrentBalance += tx.Amount
	case domain.TransactionDebit:
		if fund.CurrentBalance < tx.Amount {
			return domain.ErrInsufficientBalance
		}
		fund.CurrentBalance -= tx.Amount
	default:
		return fmt.Errorf("invalid transaction type %q", tx.Type)
	}

	if err := dbTx.Model(&fund).Update("current_balance", fund.CurrentBalance).Error; err != nil {
		return err
	}
	return dbTx.Omit(clause.Associations).Create(tx).Error
}
