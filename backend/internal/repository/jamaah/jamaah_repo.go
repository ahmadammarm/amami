package jamaah

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListFilter struct {
	Query    string
	Mustahik *bool
	MinScore *float64
	// OrderBy must be one of the fixed clauses chosen by the service, never user input.
	OrderBy string
}

type JamaahRepository interface {
	Create(jamaah *domain.Jamaah) error
	Update(jamaah *domain.Jamaah) error
	GetByID(id uuid.UUID) (*domain.Jamaah, error)
	GetPaginated(page int, limit int, filter ListFilter) ([]domain.Jamaah, int64, error)
	GetSurvey(jamaahID uuid.UUID) (*domain.MustahikData, error)
	SaveSurvey(survey *domain.MustahikData, score float64) error
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

// Update reports gorm.ErrRecordNotFound if the jamaah was deleted after it was read.
func (r *jamaahRepository) Update(jamaah *domain.Jamaah) error {
	return database.ExpectRows(r.db.Model(jamaah).
		Select("full_name", "phone", "address", "is_mustahik", "updated_at").
		Updates(jamaah))
}

func (r *jamaahRepository) GetByID(id uuid.UUID) (*domain.Jamaah, error) {
	var j domain.Jamaah
	if err := r.db.First(&j, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *jamaahRepository) GetPaginated(page int, limit int, filter ListFilter) ([]domain.Jamaah, int64, error) {
	var jamaahs []domain.Jamaah
	var total int64

	dbQuery := r.db.Model(&domain.Jamaah{})
	if filter.Query != "" {
		dbQuery = dbQuery.Where("full_name ILIKE ?", "%"+filter.Query+"%")
	}
	if filter.Mustahik != nil {
		dbQuery = dbQuery.Where("is_mustahik = ?", *filter.Mustahik)
	}
	if filter.MinScore != nil {
		dbQuery = dbQuery.Where("mustahik_score >= ?", *filter.MinScore)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := dbQuery.Order(filter.OrderBy).Offset(offset).Limit(limit).Find(&jamaahs).Error; err != nil {
		return nil, 0, err
	}

	return jamaahs, total, nil
}

func (r *jamaahRepository) GetSurvey(jamaahID uuid.UUID) (*domain.MustahikData, error) {
	var survey domain.MustahikData
	if err := r.db.First(&survey, "jamaah_id = ?", jamaahID).Error; err != nil {
		return nil, err
	}
	return &survey, nil
}

// SaveSurvey upserts the survey and stores the resulting score on the jamaah
// in one transaction, so the two never disagree. The jamaah row is locked
// first: a missing or concurrently deleted jamaah is gorm.ErrRecordNotFound.
func (r *jamaahRepository) SaveSurvey(survey *domain.MustahikData, score float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := database.LockForUpdate[domain.Jamaah](tx, survey.JamaahID); err != nil {
			return err
		}

		err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "jamaah_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"total_income", "dependents_count", "house_status", "has_disability",
				"criteria_notes", "surveyed_by", "updated_at",
			}),
		}).Create(survey).Error
		if err != nil {
			return err
		}

		return database.ExpectRows(tx.Model(&domain.Jamaah{}).Where("id = ?", survey.JamaahID).
			Update("mustahik_score", score))
	})
}
