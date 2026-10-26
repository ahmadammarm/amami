package inventory

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Asset changes whose outcome depends on the asset's current state (its loan
// status) take a callback: the repository locks the row, the service's rule
// decides against the locked state, and the write happens in the same
// transaction. Lock order is always loan before asset, so concurrent loan
// operations cannot deadlock.
type (
	// AssetEdit applies a change to a locked asset, or rejects it.
	AssetEdit func(asset *domain.Asset) error
	// LoanBuilder returns the loan for a locked asset, or rejects lending it.
	LoanBuilder func(asset *domain.Asset) (*domain.AssetLoan, error)
	// LoanReturn records a return on a locked loan, or rejects it.
	LoanReturn func(loan *domain.AssetLoan) error
)

type InventoryRepository interface {
	CreateAsset(asset *domain.Asset) error
	GetAllAssets(page int, limit int) ([]domain.Asset, int64, error)
	// UpdateAsset locks the asset, applies edit and saves the result.
	UpdateAsset(id uuid.UUID, edit AssetEdit) (*domain.Asset, error)
	// DeleteAsset locks the asset, lets check veto the delete, then soft-deletes it.
	DeleteAsset(id uuid.UUID, check AssetEdit) error

	// CreateLoan locks the asset, builds the loan, inserts it and marks the asset LOANED.
	CreateLoan(assetID uuid.UUID, build LoanBuilder) (*domain.AssetLoan, error)
	// ReturnLoan locks the loan, records the return and marks the asset GOOD.
	ReturnLoan(loanID uuid.UUID, apply LoanReturn) (*domain.AssetLoan, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) CreateAsset(asset *domain.Asset) error {
	return r.db.Create(asset).Error
}

func (r *inventoryRepository) GetAllAssets(page int, limit int) ([]domain.Asset, int64, error) {
	var assets []domain.Asset
	var total int64

	if err := r.db.Model(&domain.Asset{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("name ASC").Offset(offset).Limit(limit).Find(&assets).Error
	return assets, total, err
}

func (r *inventoryRepository) UpdateAsset(id uuid.UUID, edit AssetEdit) (*domain.Asset, error) {
	var updated *domain.Asset
	err := r.db.Transaction(func(tx *gorm.DB) error {
		asset, err := database.LockForUpdate[domain.Asset](tx, id)
		if err != nil {
			return err
		}
		if err := edit(asset); err != nil {
			return err
		}
		res := tx.Model(asset).
			Select("name", "sku", "purchase_date", "purchase_price", "current_status", "location").
			Updates(asset)
		if err := database.ExpectRows(res); err != nil {
			return err
		}
		updated = asset
		return nil
	})
	return updated, err
}

// DeleteAsset soft-deletes the asset; loan history keeps referencing it.
func (r *inventoryRepository) DeleteAsset(id uuid.UUID, check AssetEdit) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		asset, err := database.LockForUpdate[domain.Asset](tx, id)
		if err != nil {
			return err
		}
		if err := check(asset); err != nil {
			return err
		}
		return database.ExpectRows(tx.Delete(asset))
	})
}

func (r *inventoryRepository) CreateLoan(assetID uuid.UUID, build LoanBuilder) (*domain.AssetLoan, error) {
	var created *domain.AssetLoan
	err := r.db.Transaction(func(tx *gorm.DB) error {
		asset, err := database.LockForUpdate[domain.Asset](tx, assetID)
		if err != nil {
			return err
		}
		loan, err := build(asset)
		if err != nil {
			return err
		}
		if err := tx.Omit(clause.Associations).Create(loan).Error; err != nil {
			return err
		}
		if err := setAssetStatus(tx, asset.ID, domain.AssetLoaned); err != nil {
			return err
		}
		created = loan
		return nil
	})
	return created, err
}

func (r *inventoryRepository) ReturnLoan(loanID uuid.UUID, apply LoanReturn) (*domain.AssetLoan, error) {
	var returned *domain.AssetLoan
	err := r.db.Transaction(func(tx *gorm.DB) error {
		loan, err := database.LockForUpdate[domain.AssetLoan](tx, loanID)
		if err != nil {
			return err
		}
		if err := apply(loan); err != nil {
			return err
		}
		res := tx.Model(loan).Select("return_date", "condition_notes").Updates(loan)
		if err := database.ExpectRows(res); err != nil {
			return err
		}
		if err := setAssetStatus(tx, loan.AssetID, domain.AssetGood); err != nil {
			return err
		}
		returned = loan
		return nil
	})
	return returned, err
}

// setAssetStatus is only called after the loan's own lock is held (loan before
// asset); the status update takes the asset's row lock implicitly.
func setAssetStatus(tx *gorm.DB, assetID uuid.UUID, status domain.AssetStatus) error {
	return database.ExpectRows(
		tx.Model(&domain.Asset{}).Where("id = ?", assetID).Update("current_status", status),
	)
}
