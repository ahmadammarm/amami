package qurban

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanRepo "github.com/ahmadammarm/amami/backend/internal/repository/qurban"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QurbanService interface {
	CreatePackage(req qurbanDto.CreatePackageRequest) (*domain.QurbanPackage, error)
	GetAllPackages(page, limit int) ([]domain.QurbanPackage, int64, error)
	UpdatePackage(id uint, req qurbanDto.CreatePackageRequest) (*domain.QurbanPackage, error)
	DeletePackage(id uint) error

	CreateBooking(req qurbanDto.CreateBookingRequest) (*domain.QurbanBooking, error)
	GetBookings() ([]domain.QurbanBooking, error)

	CreateAnimal(req qurbanDto.CreateAnimalRequest) (*domain.QurbanAnimal, error)
	GetAnimals() ([]domain.QurbanAnimal, error)
	UpdateAnimal(id string, req qurbanDto.CreateAnimalRequest) (*domain.QurbanAnimal, error)
	DeleteAnimal(id string) error
}

type qurbanService struct {
	repo qurbanRepo.QurbanRepository
}

func NewQurbanService(repo qurbanRepo.QurbanRepository) QurbanService {
	return &qurbanService{repo}
}

func (s *qurbanService) CreatePackage(req qurbanDto.CreatePackageRequest) (*domain.QurbanPackage, error) {
	pkg := &domain.QurbanPackage{
		Name:           req.Name,
		Type:           req.Type,
		Price:          req.Price,
		YearHijri:      req.YearHijri,
		StockTotal:     req.StockTotal,
		StockRemaining: req.StockTotal,
	}

	if err := s.repo.CreatePackage(pkg); err != nil {
		return nil, err
	}
	return pkg, nil
}

func (s *qurbanService) GetAllPackages(page, limit int) ([]domain.QurbanPackage, int64, error) {
	return s.repo.GetAllPackages(page, limit)
}

func (s *qurbanService) UpdatePackage(id uint, req qurbanDto.CreatePackageRequest) (*domain.QurbanPackage, error) {
	pkg, err := s.repo.GetPackageByID(id)
	if err != nil {
		return nil, errors.New("package not found")
	}

	pkg.Name = req.Name
	pkg.Type = req.Type
	pkg.Price = req.Price
	pkg.YearHijri = req.YearHijri
	// Prevent negative stock adjustments loosely, or assume admin knows better.
	stockDiff := req.StockTotal - pkg.StockTotal
	pkg.StockTotal = req.StockTotal
	pkg.StockRemaining += stockDiff

	if pkg.StockRemaining < 0 {
		return nil, errors.New("cannot update stock total resulting in negative remaining stock")
	}

	if err := s.repo.UpdatePackage(pkg); err != nil {
		return nil, err
	}
	return pkg, nil
}

func (s *qurbanService) DeletePackage(id uint) error {
	return s.repo.DeletePackage(id)
}

func (s *qurbanService) CreateBooking(req qurbanDto.CreateBookingRequest) (*domain.QurbanBooking, error) {
	var newBooking *domain.QurbanBooking

	err := s.repo.Transaction(func(tx *gorm.DB) error {
		var pkg domain.QurbanPackage
		// Select for update to prevent concurrent booking race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pkg, "id = ?", req.PackageID).Error; err != nil {
			return errors.New("package not found")
		}

		if pkg.StockRemaining <= 0 {
			return errors.New("package out of stock")
		}

		// Decrement stock
		pkg.StockRemaining -= 1
		if err := tx.Save(&pkg).Error; err != nil {
			return err
		}

		newBooking = &domain.QurbanBooking{
			ShohibulID:    req.ShohibulID,
			PackageID:     req.PackageID,
			PaymentStatus: req.PaymentStatus,
			TotalAmount:   req.TotalAmount,
		}

		if err := tx.Create(newBooking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newBooking, nil
}

func (s *qurbanService) GetBookings() ([]domain.QurbanBooking, error) {
	return s.repo.GetBookings()
}

func (s *qurbanService) CreateAnimal(req qurbanDto.CreateAnimalRequest) (*domain.QurbanAnimal, error) {
	animal := &domain.QurbanAnimal{
		TagNumber:  req.TagNumber,
		Type:       req.Type,
		Weight:     req.Weight,
		Status:     req.Status,
		VendorInfo: req.VendorInfo,
	}

	if err := s.repo.CreateAnimal(animal); err != nil {
		return nil, err
	}
	return animal, nil
}

func (s *qurbanService) GetAnimals() ([]domain.QurbanAnimal, error) {
	return s.repo.GetAnimals()
}

func (s *qurbanService) UpdateAnimal(id string, req qurbanDto.CreateAnimalRequest) (*domain.QurbanAnimal, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid animal id")
	}

	animal, err := s.repo.GetAnimalByID(uid.String())
	if err != nil {
		return nil, errors.New("animal not found")
	}

	animal.TagNumber = req.TagNumber
	animal.Type = req.Type
	animal.Weight = req.Weight
	animal.Status = req.Status
	animal.VendorInfo = req.VendorInfo

	if err := s.repo.UpdateAnimal(animal); err != nil {
		return nil, err
	}
	return animal, nil
}

func (s *qurbanService) DeleteAnimal(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteAnimal(uid.String())
}
