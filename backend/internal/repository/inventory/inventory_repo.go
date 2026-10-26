package inventory

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryRepository interface {
	// Assets
	CreateAsset(asset *domain.Asset) error
	GetAllAssets(page int, limit int) ([]domain.Asset, int64, error)
	GetAssetByID(id uuid.UUID) (*domain.Asset, error)
	UpdateAsset(asset *domain.Asset) error
	DeleteAsset(id uuid.UUID) error

	// Loans
	GetLoanByID(id uuid.UUID) (*domain.AssetLoan, error)
	// CreateLoan inserts the loan and marks the asset LOANED atomically.
	CreateLoan(loan *domain.AssetLoan) error
	// ReturnLoan saves the returned loan and marks the asset GOOD atomically.
	ReturnLoan(loan *domain.AssetLoan) error
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

func (r *inventoryRepository) GetAssetByID(id uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *inventoryRepository) UpdateAsset(asset *domain.Asset) error {
	return r.db.Model(asset).
		Select("name", "sku", "purchase_date", "purchase_price", "current_status", "location").
		Updates(asset).Error
}

// DeleteAsset soft-deletes the asset; loan history keeps referencing it.
func (r *inventoryRepository) DeleteAsset(id uuid.UUID) error {
	res := r.db.Delete(&domain.Asset{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *inventoryRepository) GetLoanByID(id uuid.UUID) (*domain.AssetLoan, error) {
	var loan domain.AssetLoan
	if err := r.db.First(&loan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *inventoryRepository) CreateLoan(loan *domain.AssetLoan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(loan).Error; err != nil {
			return err
		}
		return setAssetStatus(tx, loan.AssetID, domain.AssetLoaned)
	})
}

func (r *inventoryRepository) ReturnLoan(loan *domain.AssetLoan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(loan).
			Select("return_date", "condition_notes").
			Updates(loan).Error
		if err != nil {
			return err
		}
		return setAssetStatus(tx, loan.AssetID, domain.AssetGood)
	})
}

func setAssetStatus(tx *gorm.DB, assetID uuid.UUID, status domain.AssetStatus) error {
	return tx.Model(&domain.Asset{}).Where("id = ?", assetID).Update("current_status", status).Error
}
