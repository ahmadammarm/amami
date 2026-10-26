package qurban

import (
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanRepo "github.com/ahmadammarm/amami/backend/internal/repository/qurban"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"gorm.io/gorm"
)

// fakeRepo mimics the locked-callback contract of UpdatePackage.
type fakeRepo struct {
	qurbanRepo.QurbanRepository
	packages map[uint]*domain.QurbanPackage
}

func (r *fakeRepo) UpdatePackage(id uint, edit qurbanRepo.PackageEdit) (*domain.QurbanPackage, error) {
	stored, ok := r.packages[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	pkg := *stored
	if err := edit(&pkg); err != nil {
		return nil, err
	}
	*stored = pkg
	return stored, nil
}

func TestUpdatePackageStock(t *testing.T) {
	// Every case starts from 10 units with 4 already booked (6 remaining).
	tests := []struct {
		name          string
		stockTotal    int
		wantRemaining int
		wantConflict  bool
	}{
		{name: "raising the quota adds remaining units", stockTotal: 15, wantRemaining: 11},
		{name: "unchanged quota", stockTotal: 10, wantRemaining: 6},
		{name: "lowering the quota removes remaining units", stockTotal: 7, wantRemaining: 3},
		{name: "quota equal to units booked", stockTotal: 4, wantRemaining: 0},
		{name: "quota below units booked", stockTotal: 3, wantConflict: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{packages: map[uint]*domain.QurbanPackage{
				1: {ID: 1, Name: "Sapi A", Type: "SAPI", Price: 20_000_000, YearHijri: 1447, StockTotal: 10, StockRemaining: 6},
			}}
			req := qurbanDto.PackageRequest{Name: "Sapi A+", Type: "SAPI", Price: 21_000_000, YearHijri: 1447, StockTotal: tt.stockTotal}

			res, err := NewQurbanService(repo).UpdatePackage(1, req)
			stored := repo.packages[1]

			if tt.wantConflict {
				if !apperror.IsKind(err, apperror.KindConflict) {
					t.Fatalf("got %v, want conflict", err)
				}
				if stored.StockTotal != 10 || stored.StockRemaining != 6 || stored.Name != "Sapi A" {
					t.Fatalf("a rejected edit must not change the package: %+v", stored)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if stored.StockTotal != tt.stockTotal || stored.StockRemaining != tt.wantRemaining {
				t.Fatalf("stock = %d/%d, want %d/%d", stored.StockRemaining, stored.StockTotal, tt.wantRemaining, tt.stockTotal)
			}
			if res.Name != "Sapi A+" || stored.Price != 21_000_000 {
				t.Fatalf("other fields not applied: %+v", stored)
			}
		})
	}

	t.Run("missing package is not found", func(t *testing.T) {
		repo := &fakeRepo{packages: map[uint]*domain.QurbanPackage{}}
		_, err := NewQurbanService(repo).UpdatePackage(9, qurbanDto.PackageRequest{StockTotal: 1})
		if !apperror.IsKind(err, apperror.KindNotFound) {
			t.Fatalf("got %v, want not found", err)
		}
	})
}
