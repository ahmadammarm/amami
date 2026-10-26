package zakat

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	financeRepo "github.com/ahmadammarm/amami/backend/internal/repository/finance"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ZakatRepository interface {
	FindJamaah(id uuid.UUID) (*domain.Jamaah, error)
	CollectZakat(donation *domain.ZakatDonation, transaction *domain.Transaction) error
	DistributeZakat(distribution *domain.ZakatDistribution, transaction *domain.Transaction) error
	GetDonations(page int, limit int) ([]domain.ZakatDonation, int64, error)
	GetDistributions(page int, limit int) ([]domain.ZakatDistribution, int64, error)
}

type zakatRepository struct {
	db *gorm.DB
}

func NewZakatRepository(db *gorm.DB) ZakatRepository {
	return &zakatRepository{db: db}
}

func (r *zakatRepository) FindJamaah(id uuid.UUID) (*domain.Jamaah, error) {
	var jamaah domain.Jamaah
	if err := r.db.First(&jamaah, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &jamaah, nil
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
// it from the ledger in the same database transaction.
func (r *zakatRepository) DistributeZakat(distribution *domain.ZakatDistribution, txModel *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if txModel != nil {
			if err := financeRepo.PostTransaction(dbTx, txModel); err != nil {
				return err
			}
			distribution.TransactionID = &txModel.ID
		}
		if err := dbTx.Omit(clause.Associations).Create(distribution).Error; err != nil {
			return err
		}
		return dbTx.First(&distribution.Mustahik, "id = ?", distribution.MustahikID).Error
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
