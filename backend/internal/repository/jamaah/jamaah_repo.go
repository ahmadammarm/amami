package jamaah

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JamaahRepository interface {
	Create(jamaah *domain.Jamaah) error
	Update(jamaah *domain.Jamaah) error
	GetByID(id uuid.UUID) (*domain.Jamaah, error)
	GetPaginated(page int, limit int, query string) ([]domain.Jamaah, int64, error)
}

type jamaahRepository struct {
	db *gorm.DB
}

func NewJamaahRepository(db *gorm.DB) JamaahRepository {
	return &jamaahRepository{db: db}
}

func (r *jamaahRepository) Create(jamaah *domain.Jamaah) error {
	return r.db.Create(jamaah).Error
}

func (r *jamaahRepository) Update(jamaah *domain.Jamaah) error {
	return r.db.Save(jamaah).Error
}

func (r *jamaahRepository) GetByID(id uuid.UUID) (*domain.Jamaah, error) {
	var j domain.Jamaah
	if err := r.db.First(&j, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *jamaahRepository) GetPaginated(page int, limit int, query string) ([]domain.Jamaah, int64, error) {
	var jamaahs []domain.Jamaah
	var total int64

	dbQuery := r.db.Model(&domain.Jamaah{})
	if query != "" {
		dbQuery = dbQuery.Where("full_name ILIKE ?", "%"+query+"%")
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := dbQuery.Order("created_at desc").Offset(offset).Limit(limit).Find(&jamaahs).Error; err != nil {
		return nil, 0, err
	}

	return jamaahs, total, nil
}
