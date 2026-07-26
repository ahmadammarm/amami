package logistics

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogisticsRepository interface {
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

	// Agenda
	CreateAgenda(agenda *domain.Agenda) error
	GetAgendas() ([]domain.Agenda, error)
	UpdateAgenda(agenda *domain.Agenda) error
}

type logisticsRepository struct {
	db *gorm.DB
}

func NewLogisticsRepository(db *gorm.DB) LogisticsRepository {
	return &logisticsRepository{db}
}

func (r *logisticsRepository) CreateAsset(asset *domain.Asset) error {
	return r.db.Create(asset).Error
}

func (r *logisticsRepository) GetAllAssets(page int, limit int) ([]domain.Asset, int64, error) {
	var assets []domain.Asset
	var total int64
	
	if err := r.db.Model(&domain.Asset{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("name ASC").Offset(offset).Limit(limit).Find(&assets).Error
	return assets, total, err
}

func (r *logisticsRepository) GetAssetByID(id uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	err := r.db.First(&asset, "id = ?", id).Error
	return &asset, err
}

func (r *logisticsRepository) UpdateAsset(asset *domain.Asset) error {
	return r.db.Save(asset).Error
}

func (r *logisticsRepository) DeleteAsset(id uuid.UUID) error {
	return r.db.Delete(&domain.Asset{}, "id = ?", id).Error
}

func (r *logisticsRepository) CreateLoan(loan *domain.AssetLoan) error {
	return r.db.Create(loan).Error
}

func (r *logisticsRepository) UpdateLoan(loan *domain.AssetLoan) error {
	return r.db.Save(loan).Error
}

func (r *logisticsRepository) GetLoanByID(id uuid.UUID) (*domain.AssetLoan, error) {
	var loan domain.AssetLoan
	err := r.db.First(&loan, "id = ?", id).Error
	return &loan, err
}

func (r *logisticsRepository) CreateAgenda(agenda *domain.Agenda) error {
	return r.db.Create(agenda).Error
}

func (r *logisticsRepository) GetAgendas() ([]domain.Agenda, error) {
	var agendas []domain.Agenda
	err := r.db.Order("start_time DESC").Find(&agendas).Error
	return agendas, err
}

func (r *logisticsRepository) UpdateAgenda(agenda *domain.Agenda) error {
	return r.db.Save(agenda).Error
}
