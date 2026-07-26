package inventory

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Assets
	CreateAsset(asset *domain.Asset) error
	GetAllAssets(page int, limit int) ([]domain.Asset, int64, error)
	GetAssetByID(id uuid.UUID) (*domain.Asset, error)
	UpdateAsset(asset *domain.Asset) error
	DeleteAsset(id uuid.UUID) error

	// Loans
	CreateLoan(loan *domain.AssetLoan) error
	UpdateLoan(loan *domain.AssetLoan) error
	GetLoanByID(id uuid.UUID) (*domain.AssetLoan, error)
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
	err := r.db.First(&asset, "id = ?", id).Error
	return &asset, err
}

func (r *inventoryRepository) UpdateAsset(asset *domain.Asset) error {
	return r.db.Save(asset).Error
}

func (r *inventoryRepository) DeleteAsset(id uuid.UUID) error {
	return r.db.Delete(&domain.Asset{}, "id = ?", id).Error
}

func (r *inventoryRepository) CreateLoan(loan *domain.AssetLoan) error {
	return r.db.Create(loan).Error
}

func (r *inventoryRepository) UpdateLoan(loan *domain.AssetLoan) error {
	return r.db.Save(loan).Error
}

func (r *inventoryRepository) GetLoanByID(id uuid.UUID) (*domain.AssetLoan, error) {
	var loan domain.AssetLoan
	err := r.db.First(&loan, "id = ?", id).Error
	return &loan, err
}
