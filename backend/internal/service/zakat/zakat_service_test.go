package zakat

import (
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	zakatDto "github.com/ahmadammarm/amami/backend/internal/dto/zakat"
	zakatRepo "github.com/ahmadammarm/amami/backend/internal/repository/zakat"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeRepo struct {
	zakatRepo.ZakatRepository
	jamaah        map[uuid.UUID]*domain.Jamaah
	collected     []*domain.Transaction
	distributed   []*domain.Transaction
	collectCalls  int
	distributeCnt int
}

func (r *fakeRepo) CollectZakat(_ *domain.ZakatDonation, tx *domain.Transaction) error {
	r.collectCalls++
	r.collected = append(r.collected, tx)
	return nil
}

// DistributeZakat mirrors the real repository: it looks up the recipient and
// runs the service's check before recording anything.
func (r *fakeRepo) DistributeZakat(dist *domain.ZakatDistribution, tx *domain.Transaction, check zakatRepo.RecipientCheck) error {
	recipient, ok := r.jamaah[dist.MustahikID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	if err := check(recipient); err != nil {
		return err
	}
	r.distributeCnt++
	r.distributed = append(r.distributed, tx)
	return nil
}

func TestCollectZakatLedgerRules(t *testing.T) {
	fundID := uint(1)
	userID := uuid.New()

	tests := []struct {
		name       string
		req        zakatDto.CollectZakatRequest
		wantErr    bool
		wantLedger bool
	}{
		{
			name:       "whole rupiah to a fund is posted",
			req:        zakatDto.CollectZakatRequest{FundID: &fundID, AmountOrQty: 250_000, Unit: domain.ZakatUnitIDR},
			wantLedger: true,
		},
		{
			name: "rice without a fund is recorded only",
			req:  zakatDto.CollectZakatRequest{AmountOrQty: 2.5, Unit: domain.ZakatUnitKG},
		},
		{
			name:    "KG posted to a fund is rejected",
			req:     zakatDto.CollectZakatRequest{FundID: &fundID, AmountOrQty: 3, Unit: domain.ZakatUnitKG},
			wantErr: true,
		},
		{
			name:    "fractional rupiah posted to a fund is rejected",
			req:     zakatDto.CollectZakatRequest{FundID: &fundID, AmountOrQty: 1000.5, Unit: domain.ZakatUnitIDR},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			_, err := NewZakatService(repo).CollectZakat(tt.req, userID)

			if tt.wantErr {
				if !apperror.IsKind(err, apperror.KindValidation) {
					t.Fatalf("got %v, want validation error", err)
				}
				if repo.collectCalls != 0 {
					t.Fatalf("nothing should be written when validation fails")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tx := repo.collected[0]
			if tt.wantLedger {
				if tx == nil || tx.Type != domain.TransactionCredit || tx.Amount != int64(tt.req.AmountOrQty) {
					t.Fatalf("unexpected ledger entry: %+v", tx)
				}
			} else if tx != nil {
				t.Fatalf("expected no ledger entry, got %+v", tx)
			}
		})
	}
}

func TestDistributeZakatRequiresMustahik(t *testing.T) {
	mustahik := &domain.Jamaah{ID: uuid.New(), IsMustahik: true}
	regular := &domain.Jamaah{ID: uuid.New(), IsMustahik: false}
	repo := &fakeRepo{jamaah: map[uuid.UUID]*domain.Jamaah{mustahik.ID: mustahik, regular.ID: regular}}
	svc := NewZakatService(repo)

	tests := []struct {
		name     string
		id       uuid.UUID
		wantKind *apperror.Kind
	}{
		{name: "registered mustahik", id: mustahik.ID},
		{name: "not a mustahik", id: regular.ID, wantKind: ptr(apperror.KindValidation)},
		{name: "unknown recipient", id: uuid.New(), wantKind: ptr(apperror.KindValidation)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := zakatDto.DistributeZakatRequest{MustahikID: tt.id, AmountOrQty: 5, Unit: domain.ZakatUnitKG}
			_, err := svc.DistributeZakat(req, uuid.New())
			if tt.wantKind == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !apperror.IsKind(err, *tt.wantKind) {
				t.Fatalf("got %v, want kind %v", err, *tt.wantKind)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
