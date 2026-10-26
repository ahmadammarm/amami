package qurban

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QurbanRepository interface {
	CreatePackage(pkg *domain.QurbanPackage) error
	GetAllPackages(page, limit int) ([]domain.QurbanPackage, int64, error)
	GetPackageByID(id uint) (*domain.QurbanPackage, error)
	UpdatePackage(pkg *domain.QurbanPackage) error
	DeletePackage(id uint) error

	// CreateBooking reserves one unit of the package and inserts the booking
	// atomically. TotalAmount is taken from the package price.
	CreateBooking(booking *domain.QurbanBooking) error
	GetBookings() ([]domain.QurbanBooking, error)

	CreateAnimal(animal *domain.QurbanAnimal) error
	GetAnimals() ([]domain.QurbanAnimal, error)
	GetAnimalByID(id uuid.UUID) (*domain.QurbanAnimal, error)
	UpdateAnimal(animal *domain.QurbanAnimal) error
	DeleteAnimal(id uuid.UUID) error
}

type qurbanRepository struct {
	db *gorm.DB
}

func NewQurbanRepository(db *gorm.DB) QurbanRepository {
	return &qurbanRepository{db}
}

func (r *qurbanRepository) CreatePackage(pkg *domain.QurbanPackage) error {
	return r.db.Create(pkg).Error
}

func (r *qurbanRepository) GetAllPackages(page, limit int) ([]domain.QurbanPackage, int64, error) {
	var pkgs []domain.QurbanPackage
	var total int64

	if err := r.db.Model(&domain.QurbanPackage{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("year_hijri DESC, name ASC").Offset(offset).Limit(limit).Find(&pkgs).Error
	return pkgs, total, err
}

func (r *qurbanRepository) GetPackageByID(id uint) (*domain.QurbanPackage, error) {
	var pkg domain.QurbanPackage
	if err := r.db.First(&pkg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *qurbanRepository) UpdatePackage(pkg *domain.QurbanPackage) error {
	return r.db.Model(pkg).
		Select("name", "type", "price", "year_hijri", "stock_total", "stock_remaining").
		Updates(pkg).Error
}

// DeletePackage soft-deletes the package; existing bookings keep referencing it.
func (r *qurbanRepository) DeletePackage(id uint) error {
	return deleteOrNotFound(r.db.Delete(&domain.QurbanPackage{}, "id = ?", id))
}

func (r *qurbanRepository) CreateBooking(booking *domain.QurbanBooking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pkg domain.QurbanPackage
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pkg, "id = ?", booking.PackageID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrPackageNotFound
		}
		if err != nil {
			return err
		}
		if pkg.StockRemaining <= 0 {
			return domain.ErrPackageOutOfStock
		}

		if err := tx.Model(&pkg).Update("stock_remaining", gorm.Expr("stock_remaining - 1")).Error; err != nil {
			return err
		}

		booking.TotalAmount = pkg.Price
		if err := tx.Omit(clause.Associations).Create(booking).Error; err != nil {
			return err
		}
		return tx.Preload("Shohibul").Preload("Package").First(booking, "id = ?", booking.ID).Error
	})
}

func (r *qurbanRepository) GetBookings() ([]domain.QurbanBooking, error) {
	var bookings []domain.QurbanBooking
	err := r.db.Preload("Shohibul", unscoped).Preload("Package", unscoped).
		Order("booking_date DESC").Find(&bookings).Error
	return bookings, err
}

func (r *qurbanRepository) CreateAnimal(animal *domain.QurbanAnimal) error {
	return r.db.Create(animal).Error
}

func (r *qurbanRepository) GetAnimals() ([]domain.QurbanAnimal, error) {
	var animals []domain.QurbanAnimal
	err := r.db.Order("tag_number ASC").Find(&animals).Error
	return animals, err
}

func (r *qurbanRepository) GetAnimalByID(id uuid.UUID) (*domain.QurbanAnimal, error) {
	var animal domain.QurbanAnimal
	if err := r.db.First(&animal, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &animal, nil
}

func (r *qurbanRepository) UpdateAnimal(animal *domain.QurbanAnimal) error {
	return r.db.Model(animal).
		Select("tag_number", "type", "weight", "status", "vendor_info").
		Updates(animal).Error
}

func (r *qurbanRepository) DeleteAnimal(id uuid.UUID) error {
	return deleteOrNotFound(r.db.Delete(&domain.QurbanAnimal{}, "id = ?", id))
}

func deleteOrNotFound(res *gorm.DB) error {
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// unscoped keeps bookings showing soft-deleted jamaah and packages.
func unscoped(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
