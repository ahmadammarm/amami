package inventory

import (
	"testing"
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventoryRepo "github.com/ahmadammarm/amami/backend/internal/repository/inventory"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fakeRepo mimics the locked-callback contract of the real repository: it
// loads the row, runs the service's rule against it and only writes when the
// rule accepts.
type fakeRepo struct {
	inventoryRepo.InventoryRepository
	assets  map[uuid.UUID]*domain.Asset
	loans   map[uuid.UUID]*domain.AssetLoan
	deleted []uuid.UUID
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{assets: map[uuid.UUID]*domain.Asset{}, loans: map[uuid.UUID]*domain.AssetLoan{}}
}

func (r *fakeRepo) addAsset(status domain.AssetStatus) *domain.Asset {
	a := &domain.Asset{ID: uuid.New(), Name: "Sound system", SKU: "SND-01", CurrentStatus: status}
	r.assets[a.ID] = a
	return a
}

func (r *fakeRepo) addLoan(loanDate time.Time, returned *time.Time) *domain.AssetLoan {
	asset := r.addAsset(domain.AssetLoaned)
	l := &domain.AssetLoan{ID: uuid.New(), AssetID: asset.ID, LoanDate: loanDate, ReturnDate: returned, ConditionNotes: "ok"}
	r.loans[l.ID] = l
	return l
}

func (r *fakeRepo) UpdateAsset(id uuid.UUID, edit inventoryRepo.AssetEdit) (*domain.Asset, error) {
	stored, ok := r.assets[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	asset := *stored
	if err := edit(&asset); err != nil {
		return nil, err
	}
	*stored = asset
	return stored, nil
}

func (r *fakeRepo) DeleteAsset(id uuid.UUID, check inventoryRepo.AssetEdit) error {
	asset, ok := r.assets[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	if err := check(asset); err != nil {
		return err
	}
	delete(r.assets, id)
	r.deleted = append(r.deleted, id)
	return nil
}

func (r *fakeRepo) CreateLoan(assetID uuid.UUID, build inventoryRepo.LoanBuilder) (*domain.AssetLoan, error) {
	asset, ok := r.assets[assetID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	loan, err := build(asset)
	if err != nil {
		return nil, err
	}
	loan.ID = uuid.New()
	r.loans[loan.ID] = loan
	asset.CurrentStatus = domain.AssetLoaned
	return loan, nil
}

func (r *fakeRepo) ReturnLoan(loanID uuid.UUID, apply inventoryRepo.LoanReturn) (*domain.AssetLoan, error) {
	stored, ok := r.loans[loanID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	loan := *stored
	if err := apply(&loan); err != nil {
		return nil, err
	}
	*stored = loan
	if asset, ok := r.assets[loan.AssetID]; ok {
		asset.CurrentStatus = domain.AssetGood
	}
	return stored, nil
}

func assertKind(t *testing.T, err error, want *apperror.Kind) {
	t.Helper()
	if want == nil {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	if !apperror.IsKind(err, *want) {
		t.Fatalf("got %v, want kind %v", err, *want)
	}
}

func ptr[T any](v T) *T { return &v }

func TestUpdateAssetKeepsLoanStatusOutOfReach(t *testing.T) {
	tests := []struct {
		name     string
		current  domain.AssetStatus
		next     domain.AssetStatus
		wantKind *apperror.Kind
	}{
		{name: "good to needs repair", current: domain.AssetGood, next: domain.AssetNeedsRepair},
		{name: "loaned stays loaned while other fields change", current: domain.AssetLoaned, next: domain.AssetLoaned},
		{name: "cannot mark as loaned", current: domain.AssetGood, next: domain.AssetLoaned, wantKind: ptr(apperror.KindConflict)},
		{name: "cannot clear loaned", current: domain.AssetLoaned, next: domain.AssetGood, wantKind: ptr(apperror.KindConflict)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			asset := repo.addAsset(tt.current)
			req := inventoryDto.AssetRequest{Name: "Renamed", SKU: asset.SKU, CurrentStatus: tt.next}

			_, err := NewInventoryService(repo).UpdateAsset(asset.ID, req)
			assertKind(t, err, tt.wantKind)

			wantName := "Renamed"
			if tt.wantKind != nil {
				wantName = "Sound system"
			}
			if got := repo.assets[asset.ID].Name; got != wantName {
				t.Fatalf("name = %q, want %q", got, wantName)
			}
		})
	}

	t.Run("missing asset is not found", func(t *testing.T) {
		_, err := NewInventoryService(newFakeRepo()).UpdateAsset(uuid.New(), inventoryDto.AssetRequest{CurrentStatus: domain.AssetGood})
		assertKind(t, err, ptr(apperror.KindNotFound))
	})
}

func TestDeleteAssetRejectsAssetsOnLoan(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.AssetStatus
		wantKind *apperror.Kind
	}{
		{name: "asset in storage", status: domain.AssetGood},
		{name: "broken asset", status: domain.AssetBroken},
		{name: "asset on loan", status: domain.AssetLoaned, wantKind: ptr(apperror.KindConflict)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			asset := repo.addAsset(tt.status)

			err := NewInventoryService(repo).DeleteAsset(asset.ID)
			assertKind(t, err, tt.wantKind)
			if wantDeleted := tt.wantKind == nil; (len(repo.deleted) == 1) != wantDeleted {
				t.Fatalf("deleted = %v, want deleted %v", repo.deleted, wantDeleted)
			}
		})
	}
}

func TestCreateLoan(t *testing.T) {
	day := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	before := day.AddDate(0, 0, -1)

	tests := []struct {
		name     string
		status   domain.AssetStatus
		dueDate  *time.Time
		wantKind *apperror.Kind
	}{
		{name: "asset in good condition", status: domain.AssetGood},
		{name: "asset already on loan", status: domain.AssetLoaned, wantKind: ptr(apperror.KindConflict)},
		{name: "asset under repair", status: domain.AssetRepair, wantKind: ptr(apperror.KindConflict)},
		{name: "due date before loan date", status: domain.AssetGood, dueDate: &before, wantKind: ptr(apperror.KindValidation)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			asset := repo.addAsset(tt.status)
			req := inventoryDto.CreateAssetLoanRequest{JamaahID: uuid.New(), LoanDate: day, DueDate: tt.dueDate}

			loan, err := NewInventoryService(repo).CreateLoan(asset.ID, req)
			assertKind(t, err, tt.wantKind)
			if tt.wantKind != nil {
				if len(repo.loans) != 0 {
					t.Fatalf("no loan should be stored when the request is rejected")
				}
				return
			}
			if loan.AssetID != asset.ID || asset.CurrentStatus != domain.AssetLoaned {
				t.Fatalf("loan %+v, asset status %s", loan, asset.CurrentStatus)
			}
		})
	}

	t.Run("borrower that does not exist", func(t *testing.T) {
		repo := newFakeRepo()
		asset := repo.addAsset(domain.AssetGood)
		failing := &fkRepo{fakeRepo: repo}
		_, err := NewInventoryService(failing).CreateLoan(asset.ID, inventoryDto.CreateAssetLoanRequest{LoanDate: day})
		assertKind(t, err, ptr(apperror.KindValidation))
	})
}

// fkRepo fails the loan insert the way PostgreSQL does for an unknown jamaah.
type fkRepo struct{ *fakeRepo }

func (r *fkRepo) CreateLoan(assetID uuid.UUID, build inventoryRepo.LoanBuilder) (*domain.AssetLoan, error) {
	if _, err := build(r.assets[assetID]); err != nil {
		return nil, err
	}
	return nil, gorm.ErrForeignKeyViolated
}

func TestReturnLoan(t *testing.T) {
	loanDate := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	returnedEarlier := loanDate.AddDate(0, 0, 2)

	tests := []struct {
		name       string
		returned   *time.Time
		returnDate time.Time
		wantKind   *apperror.Kind
	}{
		{name: "open loan", returnDate: loanDate.AddDate(0, 0, 3)},
		{name: "same day return", returnDate: loanDate},
		{name: "already returned", returned: &returnedEarlier, returnDate: loanDate.AddDate(0, 0, 3), wantKind: ptr(apperror.KindConflict)},
		{name: "return before loan date", returnDate: loanDate.AddDate(0, 0, -1), wantKind: ptr(apperror.KindValidation)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			loan := repo.addLoan(loanDate, tt.returned)
			req := inventoryDto.ReturnAssetLoanRequest{ReturnDate: tt.returnDate, ConditionNotes: "scratched"}

			res, err := NewInventoryService(repo).ReturnLoan(loan.ID, req)
			assertKind(t, err, tt.wantKind)
			if tt.wantKind != nil {
				if repo.loans[loan.ID].ConditionNotes != "ok" {
					t.Fatalf("a rejected return must not change the loan")
				}
				return
			}
			if res.ReturnDate == nil || !res.ReturnDate.Equal(tt.returnDate) {
				t.Fatalf("return date = %v, want %v", res.ReturnDate, tt.returnDate)
			}
			if res.ConditionNotes != "ok | Return: scratched" {
				t.Fatalf("condition notes = %q", res.ConditionNotes)
			}
			if repo.assets[loan.AssetID].CurrentStatus != domain.AssetGood {
				t.Fatalf("asset should be back in GOOD condition")
			}
		})
	}

	t.Run("missing loan is not found", func(t *testing.T) {
		_, err := NewInventoryService(newFakeRepo()).ReturnLoan(uuid.New(), inventoryDto.ReturnAssetLoanRequest{ReturnDate: loanDate})
		assertKind(t, err, ptr(apperror.KindNotFound))
	})
}
