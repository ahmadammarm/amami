package finance

import (
	"errors"
	"fmt"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostTransaction appends tx to the ledger and applies it to the fund balance.
// It must run inside a database transaction: the fund row is locked so
// concurrent postings to the same fund are serialized. Every ledger write in
// the system (manual entries, reversals, zakat) goes through this function.
func PostTransaction(dbTx *gorm.DB, tx *domain.Transaction) error {
	fund, err := database.LockForUpdate[domain.Fund](dbTx, tx.FundID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrFundNotFound
	}
	if err != nil {
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

	if err := dbTx.Model(fund).Update("current_balance", fund.CurrentBalance).Error; err != nil {
		return err
	}
	return dbTx.Omit(clause.Associations).Create(tx).Error
}
