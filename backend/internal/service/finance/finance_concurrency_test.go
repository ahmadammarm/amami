package finance

import (
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	financeDto "github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/ahmadammarm/amami/backend/internal/testutil/pgtest"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
)

// Concurrent debits against one fund must never overdraw it, and the stored
// balance must always equal the sum of the ledger.
func TestConcurrentDebitsNeverOverdraw(t *testing.T) {
	const debits, amount, covered = 20, 100_000, 5

	db := pgtest.New(t)
	svc := NewFinanceService(financeRepo.NewFinanceRepository(db))
	user := pgtest.CreateUser(t, db)

	fund := &domain.Fund{Name: "Kas Operasional", Code: "OPS"}
	if err := db.Create(fund).Error; err != nil {
		t.Fatalf("create fund: %v", err)
	}
	opening := financeDto.CreateTransactionRequest{FundID: fund.ID, Type: domain.TransactionCredit, Amount: amount * covered, Category: "INFAQ"}
	if _, err := svc.CreateTransaction(opening, user.ID); err != nil {
		t.Fatalf("opening credit: %v", err)
	}

	debit := financeDto.CreateTransactionRequest{FundID: fund.ID, Type: domain.TransactionDebit, Amount: amount, Category: "OPERASIONAL"}
	errs := pgtest.RunConcurrently(debits, func(int) error {
		_, err := svc.CreateTransaction(debit, user.ID)
		return err
	})

	if got := pgtest.CountNil(errs); got != covered {
		t.Fatalf("%d debits succeeded, want %d", got, covered)
	}
	for _, err := range errs {
		if err != nil && !apperror.IsKind(err, apperror.KindConflict) {
			t.Fatalf("rejected debits must be conflicts, got %v", err)
		}
	}

	var balance, ledger int64
	if err := db.Model(&domain.Fund{}).Where("id = ?", fund.ID).Pluck("current_balance", &balance).Error; err != nil {
		t.Fatalf("read balance: %v", err)
	}
	err := db.Model(&domain.Transaction{}).Where("fund_id = ?", fund.ID).
		Select("COALESCE(SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE -amount END), 0)").
		Scan(&ledger).Error
	if err != nil {
		t.Fatalf("sum ledger: %v", err)
	}
	if balance != 0 || balance != ledger {
		t.Fatalf("balance = %d, ledger sum = %d, want both 0", balance, ledger)
	}
}
