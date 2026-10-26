package inventory

import (
	"testing"
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventoryRepo "github.com/ahmadammarm/amami/backend/internal/repository/inventory"
	"github.com/ahmadammarm/amami/backend/internal/testutil/pgtest"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
)

const contenders = 20

func newDBService(t *testing.T) (InventoryService, *domain.Jamaah, *inventoryDto.AssetResponse) {
	t.Helper()
	db := pgtest.New(t)
	svc := NewInventoryService(inventoryRepo.NewInventoryRepository(db))
	borrower := pgtest.CreateJamaah(t, db, false)
	asset, err := svc.CreateAsset(inventoryDto.AssetRequest{
		Name: "Tenda", SKU: "TND-" + uuid.NewString()[:8], CurrentStatus: domain.AssetGood,
	})
	if err != nil {
		t.Fatalf("create asset: %v", err)
	}
	return svc, borrower, asset
}

func requireConflicts(t *testing.T, errs []error) {
	t.Helper()
	for _, err := range errs {
		if err != nil && !apperror.IsKind(err, apperror.KindConflict) {
			t.Fatalf("losing requests must fail with a conflict, got %v", err)
		}
	}
}

func TestConcurrentLoansOfOneAsset(t *testing.T) {
	svc, borrower, asset := newDBService(t)
	req := inventoryDto.CreateAssetLoanRequest{JamaahID: borrower.ID, LoanDate: time.Now()}

	errs := pgtest.RunConcurrently(contenders, func(int) error {
		_, err := svc.CreateLoan(asset.ID, req)
		return err
	})

	if got := pgtest.CountNil(errs); got != 1 {
		t.Fatalf("%d of %d loans succeeded, want exactly 1: %v", got, contenders, errs)
	}
	requireConflicts(t, errs)
}

func TestConcurrentReturnsOfOneLoan(t *testing.T) {
	svc, borrower, asset := newDBService(t)
	loanDate := time.Now().Add(-time.Hour)
	loan, err := svc.CreateLoan(asset.ID, inventoryDto.CreateAssetLoanRequest{JamaahID: borrower.ID, LoanDate: loanDate})
	if err != nil {
		t.Fatalf("create loan: %v", err)
	}

	errs := pgtest.RunConcurrently(contenders, func(i int) error {
		_, err := svc.ReturnLoan(loan.ID, inventoryDto.ReturnAssetLoanRequest{ReturnDate: time.Now(), ConditionNotes: "ok"})
		return err
	})

	if got := pgtest.CountNil(errs); got != 1 {
		t.Fatalf("%d of %d returns succeeded, want exactly 1: %v", got, contenders, errs)
	}
	requireConflicts(t, errs)

	// The asset is lendable again, exactly once.
	_, err = svc.CreateLoan(asset.ID, inventoryDto.CreateAssetLoanRequest{JamaahID: borrower.ID, LoanDate: time.Now()})
	if err != nil {
		t.Fatalf("asset should be lendable after the return: %v", err)
	}
}

func TestEditDuringLoanKeepsLoanedStatus(t *testing.T) {
	svc, borrower, asset := newDBService(t)
	edit := inventoryDto.AssetRequest{Name: "Tenda besar", SKU: asset.SKU, CurrentStatus: domain.AssetNeedsRepair}
	loan := inventoryDto.CreateAssetLoanRequest{JamaahID: borrower.ID, LoanDate: time.Now()}

	// Half the goroutines lend the asset, half edit it. Whatever the order,
	// an edit must never overwrite LOANED set by a loan.
	errs := pgtest.RunConcurrently(contenders, func(i int) error {
		if i%2 == 0 {
			_, err := svc.CreateLoan(asset.ID, loan)
			return err
		}
		_, err := svc.UpdateAsset(asset.ID, edit)
		return err
	})
	requireConflicts(t, errs)

	loans := 0
	for i := 0; i < contenders; i += 2 {
		if errs[i] == nil {
			loans++
		}
	}
	if loans > 1 {
		t.Fatalf("%d loans of one asset succeeded", loans)
	}

	page, err := svc.GetAllAssets(1, 10)
	if err != nil {
		t.Fatalf("list assets: %v", err)
	}
	status := page.Items[0].CurrentStatus
	if (loans == 1) != (status == domain.AssetLoaned) {
		t.Fatalf("asset status %s does not match %d active loans", status, loans)
	}
}
