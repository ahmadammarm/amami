package zakat

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecipientCheck accepts or rejects the locked recipient of a distribution.
type RecipientCheck func(recipient *domain.Jamaah) error

type ZakatRepository interface {
	CollectZakat(donation *domain.ZakatDonation, transaction *domain.Transaction) error
	// DistributeZakat share-locks the recipient, lets check reject them, then
	// records the distribution. A missing recipient is gorm.ErrRecordNotFound.
	DistributeZakat(distribution *domain.ZakatDistribution, transaction *domain.Transaction, check RecipientCheck) error
	GetDonations(page int, limit int) ([]domain.ZakatDonation, int64, error)
	GetDistributions(page int, limit int) ([]domain.ZakatDistribution, int64, error)
}

type zakatRepository struct {
	db *gorm.DB
}

func NewZakatRepository(db *gorm.DB) ZakatRepository {
	return &zakatRepository{db: db}
}

// CollectZakat records a donation and, when txModel is set, posts it to the
// ledger in the same database transaction.
func (r *zakatRepository) CollectZakat(donation *domain.ZakatDonation, txModel *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if txModel != nil {
			if err := financeRepo.PostTransaction(dbTx, txModel); err != nil {
				return err
			}
			donation.TransactionID = &txModel.ID
		}
		if err := dbTx.Omit(clause.Associations).Create(donation).Error; err != nil {
			return err
		}
		return dbTx.First(&donation.Muzakki, "id = ?", donation.MuzakkiID).Error
	})
}

// DistributeZakat records a distribution and, when txModel is set, withdraws
// it from the ledger in the same database transaction. The recipient row is
// share-locked first (lock order: jamaah, then fund), so their mustahik status
// cannot change or the row be deleted until the distribution commits.
func (r *zakatRepository) DistributeZakat(distribution *domain.ZakatDistribution, txModel *domain.Transaction, check RecipientCheck) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		recipient, err := database.LockForShare[domain.Jamaah](dbTx, distribution.MustahikID)
		if err != nil {
			return err
		}
		if err := check(recipient); err != nil {
			return err
		}

		if txModel != nil {
			if err := financeRepo.PostTransaction(dbTx, txModel); err != nil {
				return err
			}
			distribution.TransactionID = &txModel.ID
		}
		if err := dbTx.Omit(clause.Associations).Create(distribution).Error; err != nil {
			return err
		}
		distribution.Mustahik = *recipient
		return nil
	})
}

func (r *zakatRepository) GetDonations(page int, limit int) ([]domain.ZakatDonation, int64, error) {
	var donations []domain.ZakatDonation
	var total int64

	if err := r.db.Model(&domain.ZakatDonation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Muzakki", unscoped).
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&donations).Error
	if err != nil {
		return nil, 0, err
	}
	return donations, total, nil
}

func (r *zakatRepository) GetDistributions(page int, limit int) ([]domain.ZakatDistribution, int64, error) {
	var distributions []domain.ZakatDistribution
	var total int64

	if err := r.db.Model(&domain.ZakatDistribution{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Mustahik", unscoped).
		Order("distributed_at desc").
		Offset(offset).Limit(limit).
		Find(&distributions).Error
	if err != nil {
		return nil, 0, err
	}
	return distributions, total, nil
}

// unscoped keeps historical zakat records showing the names of soft-deleted jamaah.
func unscoped(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
