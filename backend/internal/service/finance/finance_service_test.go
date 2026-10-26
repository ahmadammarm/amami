package finance

import (
	"errors"
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	financeDto "github.com/ahmadammarm/amami/backend/internal/dto/finance"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fakeRepo mimics the reversal contract of the real repository: it looks up the
// original, rejects it if already reversed, then posts what the builder returns.
type fakeRepo struct {
	financeRepo.FinanceRepository
	transactions map[uuid.UUID]*domain.Transaction
	reversed     map[uuid.UUID]bool
	posted       []*domain.Transaction
	postErr      error
}

func newFakeRepo(txs ...*domain.Transaction) *fakeRepo {
	r := &fakeRepo{transactions: map[uuid.UUID]*domain.Transaction{}, reversed: map[uuid.UUID]bool{}}
	for _, tx := range txs {
		r.transactions[tx.ID] = tx
	}
	return r
}

func (r *fakeRepo) ReverseTransaction(originalID uuid.UUID, build financeRepo.ReversalBuilder) (*domain.Transaction, error) {
	original, ok := r.transactions[originalID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	if r.reversed[originalID] {
		return nil, domain.ErrAlreadyReversed
	}
	reversal, err := build(original)
	if err != nil {
		return nil, err
	}
	if r.postErr != nil {
		return nil, r.postErr
	}
	r.reversed[originalID] = true
	r.posted = append(r.posted, reversal)
	return reversal, nil
}

func TestReverseTransaction(t *testing.T) {
	userID := uuid.New()
	credit := &domain.Transaction{ID: uuid.New(), FundID: 7, Type: domain.TransactionCredit, Amount: 150_000, Category: "INFAQ"}
	reversalOf := credit.ID
	reversalEntry := &domain.Transaction{ID: uuid.New(), FundID: 7, Type: domain.TransactionDebit, Amount: 150_000, ReversalOfID: &reversalOf}
	req := financeDto.ReverseTransactionRequest{Reason: "double entry"}

	t.Run("posts an opposite entry linked to the original", func(t *testing.T) {
		repo := newFakeRepo(credit)
		svc := NewFinanceService(repo)

		res, err := svc.ReverseTransaction(credit.ID, req, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.posted) != 1 {
			t.Fatalf("expected one posted reversal, got %d", len(repo.posted))
		}
		posted := repo.posted[0]
		if posted.Type != domain.TransactionDebit || posted.Amount != credit.Amount || posted.FundID != credit.FundID {
			t.Fatalf("reversal does not mirror original: %+v", posted)
		}
		if posted.ReversalOfID == nil || *posted.ReversalOfID != credit.ID {
			t.Fatalf("reversal is not linked to original")
		}
		if posted.Category != domain.CategoryReversal || posted.Description != req.Reason || posted.CreatedBy != userID {
			t.Fatalf("unexpected reversal metadata: %+v", posted)
		}
		if res.ReversalOfID == nil || *res.ReversalOfID != credit.ID {
			t.Fatalf("response is missing reversal_of_id")
		}
	})

	tests := []struct {
		name     string
		repo     *fakeRepo
		id       uuid.UUID
		wantKind apperror.Kind
	}{
		{name: "unknown transaction", repo: newFakeRepo(), id: uuid.New(), wantKind: apperror.KindNotFound},
		{name: "already reversed", repo: func() *fakeRepo {
			r := newFakeRepo(credit)
			r.reversed[credit.ID] = true
			return r
		}(), id: credit.ID, wantKind: apperror.KindConflict},
		{name: "reversing a reversal", repo: newFakeRepo(reversalEntry), id: reversalEntry.ID, wantKind: apperror.KindConflict},
		{name: "insufficient balance", repo: func() *fakeRepo {
			r := newFakeRepo(credit)
			r.postErr = domain.ErrInsufficientBalance
			return r
		}(), id: credit.ID, wantKind: apperror.KindConflict},
		{name: "unexpected failure", repo: func() *fakeRepo {
			r := newFakeRepo(credit)
			r.postErr = errors.New("connection reset")
			return r
		}(), id: credit.ID, wantKind: apperror.KindInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFinanceService(tt.repo).ReverseTransaction(tt.id, req, userID)
			if !apperror.IsKind(err, tt.wantKind) {
				t.Fatalf("got error %v, want kind %v", err, tt.wantKind)
			}
			if len(tt.repo.posted) != 0 {
				t.Fatalf("nothing should be posted on failure")
			}
		})
	}
}

func TestBuildReversalFlipsDirection(t *testing.T) {
	for _, txType := range []domain.TransactionType{domain.TransactionCredit, domain.TransactionDebit} {
		original := &domain.Transaction{ID: uuid.New(), Type: txType, Amount: 1}
		reversal, err := BuildReversal(original, "fix", uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reversal.Type == txType {
			t.Fatalf("reversal of %s kept the same type", txType)
		}
	}
}
