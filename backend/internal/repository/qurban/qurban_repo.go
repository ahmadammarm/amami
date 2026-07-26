package qurban

import (
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type QurbanRepository interface {
	CreatePackage(pkg *domain.QurbanPackage) error
	GetAllPackages(page, limit int) ([]domain.QurbanPackage, int64, error)
	GetPackageByID(id uint) (*domain.QurbanPackage, error)
	UpdatePackage(pkg *domain.QurbanPackage) error
	DeletePackage(id uint) error

	CreateBooking(booking *domain.QurbanBooking) error
	GetBookings() ([]domain.QurbanBooking, error)

	CreateAnimal(animal *domain.QurbanAnimal) error
	GetAnimals() ([]domain.QurbanAnimal, error)
	GetAnimalByID(id string) (*domain.QurbanAnimal, error)
	UpdateAnimal(animal *domain.QurbanAnimal) error
	DeleteAnimal(id string) error

	Transaction(fn func(tx *gorm.DB) error) error
}

type qurbanRepository struct {
	db *gorm.DB
}

func NewQurbanRepository(db *gorm.DB) QurbanRepository {
	return &qurbanRepository{db}
}

func (r *qurbanRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
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
	err := r.db.First(&pkg, "id = ?", id).Error
	return &pkg, err
}

func (r *qurbanRepository) UpdatePackage(pkg *domain.QurbanPackage) error {
	return r.db.Save(pkg).Error
}

func (r *qurbanRepository) DeletePackage(id uint) error {
	return r.db.Delete(&domain.QurbanPackage{}, "id = ?", id).Error
}

func (r *qurbanRepository) CreateBooking(booking *domain.QurbanBooking) error {
	return r.db.Create(booking).Error
}

func (r *qurbanRepository) GetBookings() ([]domain.QurbanBooking, error) {
	var bookings []domain.QurbanBooking
	err := r.db.Preload("Shohibul").Preload("Package").Order("booking_date DESC").Find(&bookings).Error
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

func (r *qurbanRepository) GetAnimalByID(id string) (*domain.QurbanAnimal, error) {
	var animal domain.QurbanAnimal
	err := r.db.First(&animal, "id = ?", id).Error
	return &animal, err
}

func (r *qurbanRepository) UpdateAnimal(animal *domain.QurbanAnimal) error {
	return r.db.Save(animal).Error
}

func (r *qurbanRepository) DeleteAnimal(id string) error {
	return r.db.Delete(&domain.QurbanAnimal{}, "id = ?", id).Error
}
