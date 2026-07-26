package zakat

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ZakatRepository interface {
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

// CollectZakat handles an atomic database transaction to update the fund balance, insert a transaction, and record the zakat donation.
func (r *zakatRepository) CollectZakat(donation *domain.ZakatDonation, txModel *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if txModel != nil {
			// 1. Lock Fund
			var fund domain.Fund
			if err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fund, txModel.FundID).Error; err != nil {
				return err
			}

			// 2. Update Balance
			fund.CurrentBalance += txModel.Amount
			if err := dbTx.Save(&fund).Error; err != nil {
				return err
			}

			// 3. Create Transaction Log
			if err := dbTx.Create(txModel).Error; err != nil {
				return err
			}
			donation.TransactionID = &txModel.ID
		}

		// 4. Create Zakat Donation
		if err := dbTx.Create(donation).Error; err != nil {
			return err
		}

		return nil
	})
}

// DistributeZakat handles atomic withdrawal from fund and records distribution.
func (r *zakatRepository) DistributeZakat(distribution *domain.ZakatDistribution, txModel *domain.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if txModel != nil {
			// 1. Lock Fund
			var fund domain.Fund
			if err := dbTx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fund, txModel.FundID).Error; err != nil {
				return err
			}

			// 2. Check Balance
			if fund.CurrentBalance < txModel.Amount {
				return errors.New("insufficient fund balance")
			}

			// 3. Update Balance
			fund.CurrentBalance -= txModel.Amount
			if err := dbTx.Save(&fund).Error; err != nil {
				return err
			}

			// 4. Create Transaction Log
			if err := dbTx.Create(txModel).Error; err != nil {
				return err
			}
			distribution.TransactionID = &txModel.ID
		}

		// 5. Create Distribution
		if err := dbTx.Create(distribution).Error; err != nil {
			return err
		}

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
	if err := r.db.Preload("Muzakki").Preload("Transaction").Joins("LEFT JOIN transactions ON transactions.id = zakat_donations.transaction_id").Order("COALESCE(transactions.created_at, zakat_donations.created_at) desc").Offset(offset).Limit(limit).Find(&donations).Error; err != nil {
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
	if err := r.db.Preload("Mustahik").Preload("Transaction").Order("distributed_at desc").Offset(offset).Limit(limit).Find(&distributions).Error; err != nil {
		return nil, 0, err
	}

	return distributions, total, nil
}
