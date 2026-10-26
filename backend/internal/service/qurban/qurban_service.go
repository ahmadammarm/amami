package qurban

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	qurbanDto "github.com/ahmadammarm/amami/backend/internal/dto/qurban"
	qurbanRepo "github.com/ahmadammarm/amami/backend/internal/repository/qurban"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QurbanService interface {
	CreatePackage(req qurbanDto.PackageRequest) (*qurbanDto.PackageResponse, error)
	GetAllPackages(page, limit int) (utils.Page[qurbanDto.PackageResponse], error)
	UpdatePackage(id uint, req qurbanDto.PackageRequest) (*qurbanDto.PackageResponse, error)
	DeletePackage(id uint) error

	CreateBooking(req qurbanDto.CreateBookingRequest) (*qurbanDto.BookingResponse, error)
	GetBookings() ([]qurbanDto.BookingResponse, error)

	CreateAnimal(req qurbanDto.AnimalRequest) (*qurbanDto.AnimalResponse, error)
	GetAnimals() ([]qurbanDto.AnimalResponse, error)
	UpdateAnimal(id uuid.UUID, req qurbanDto.AnimalRequest) (*qurbanDto.AnimalResponse, error)
	DeleteAnimal(id uuid.UUID) error
}

type qurbanService struct {
	repo qurbanRepo.QurbanRepository
}

func NewQurbanService(repo qurbanRepo.QurbanRepository) QurbanService {
	return &qurbanService{repo}
}

func (s *qurbanService) CreatePackage(req qurbanDto.PackageRequest) (*qurbanDto.PackageResponse, error) {
	pkg := &domain.QurbanPackage{
		Name:           req.Name,
		Type:           req.Type,
		Price:          req.Price,
		YearHijri:      req.YearHijri,
		StockTotal:     req.StockTotal,
		StockRemaining: req.StockTotal,
	}

	if err := s.repo.CreatePackage(pkg); err != nil {
		return nil, apperror.Internal("failed to create package", err)
	}

	res := qurbanDto.ToPackageResponse(pkg)
	return &res, nil
}

func (s *qurbanService) GetAllPackages(page, limit int) (utils.Page[qurbanDto.PackageResponse], error) {
	pkgs, total, err := s.repo.GetAllPackages(page, limit)
	if err != nil {
		return utils.Page[qurbanDto.PackageResponse]{}, apperror.Internal("failed to fetch packages", err)
	}

	res := make([]qurbanDto.PackageResponse, 0, len(pkgs))
	for i := range pkgs {
		res = append(res, qurbanDto.ToPackageResponse(&pkgs[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}

func (s *qurbanService) UpdatePackage(id uint, req qurbanDto.PackageRequest) (*qurbanDto.PackageResponse, error) {
	pkg, err := s.repo.GetPackageByID(id)
	if err != nil {
		return nil, notFoundOr(err, "package not found", "failed to fetch package")
	}

	// Changing the total shifts the remaining stock by the same amount, so
	// units already booked stay booked.
	remaining := pkg.StockRemaining + (req.StockTotal - pkg.StockTotal)
	if remaining < 0 {
		return nil, apperror.Conflict("stock total cannot be lower than the number of units already booked")
	}

	pkg.Name = req.Name
	pkg.Type = req.Type
	pkg.Price = req.Price
	pkg.YearHijri = req.YearHijri
	pkg.StockTotal = req.StockTotal
	pkg.StockRemaining = remaining

	if err := s.repo.UpdatePackage(pkg); err != nil {
		return nil, apperror.Internal("failed to update package", err)
	}

	res := qurbanDto.ToPackageResponse(pkg)
	return &res, nil
}

func (s *qurbanService) DeletePackage(id uint) error {
	if err := s.repo.DeletePackage(id); err != nil {
		return notFoundOr(err, "package not found", "failed to delete package")
	}
	return nil
}

func (s *qurbanService) CreateBooking(req qurbanDto.CreateBookingRequest) (*qurbanDto.BookingResponse, error) {
	booking := &domain.QurbanBooking{
		ShohibulID:    req.ShohibulID,
		PackageID:     req.PackageID,
		PaymentStatus: req.PaymentStatus,
	}

	if err := s.repo.CreateBooking(booking); err != nil {
		switch {
		case errors.Is(err, domain.ErrPackageNotFound):
			return nil, apperror.Validation(err.Error())
		case errors.Is(err, domain.ErrPackageOutOfStock):
			return nil, apperror.Conflict(err.Error())
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return nil, apperror.Validation("shohibul not found")
		default:
			return nil, apperror.Internal("failed to create booking", err)
		}
	}

	res := qurbanDto.ToBookingResponse(booking)
	return &res, nil
}

func (s *qurbanService) GetBookings() ([]qurbanDto.BookingResponse, error) {
	bookings, err := s.repo.GetBookings()
	if err != nil {
		return nil, apperror.Internal("failed to fetch bookings", err)
	}

	res := make([]qurbanDto.BookingResponse, 0, len(bookings))
	for i := range bookings {
		res = append(res, qurbanDto.ToBookingResponse(&bookings[i]))
	}
	return res, nil
}

func (s *qurbanService) CreateAnimal(req qurbanDto.AnimalRequest) (*qurbanDto.AnimalResponse, error) {
	animal := &domain.QurbanAnimal{}
	applyAnimalRequest(animal, req)

	if err := s.repo.CreateAnimal(animal); err != nil {
		return nil, mapAnimalWriteError(err, "failed to create animal")
	}

	res := qurbanDto.ToAnimalResponse(animal)
	return &res, nil
}

func (s *qurbanService) GetAnimals() ([]qurbanDto.AnimalResponse, error) {
	animals, err := s.repo.GetAnimals()
	if err != nil {
		return nil, apperror.Internal("failed to fetch animals", err)
	}

	res := make([]qurbanDto.AnimalResponse, 0, len(animals))
	for i := range animals {
		res = append(res, qurbanDto.ToAnimalResponse(&animals[i]))
	}
	return res, nil
}

func (s *qurbanService) UpdateAnimal(id uuid.UUID, req qurbanDto.AnimalRequest) (*qurbanDto.AnimalResponse, error) {
	animal, err := s.repo.GetAnimalByID(id)
	if err != nil {
		return nil, notFoundOr(err, "animal not found", "failed to fetch animal")
	}

	applyAnimalRequest(animal, req)
	if err := s.repo.UpdateAnimal(animal); err != nil {
		return nil, mapAnimalWriteError(err, "failed to update animal")
	}

	res := qurbanDto.ToAnimalResponse(animal)
	return &res, nil
}

func (s *qurbanService) DeleteAnimal(id uuid.UUID) error {
	if err := s.repo.DeleteAnimal(id); err != nil {
		return notFoundOr(err, "animal not found", "failed to delete animal")
	}
	return nil
}

func applyAnimalRequest(animal *domain.QurbanAnimal, req qurbanDto.AnimalRequest) {
	animal.TagNumber = req.TagNumber
	animal.Type = req.Type
	animal.Weight = req.Weight
	animal.Status = req.Status
	animal.VendorInfo = req.VendorInfo
}

func mapAnimalWriteError(err error, msg string) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.Conflict("an animal with this tag number already exists")
	}
	return apperror.Internal(msg, err)
}

func notFoundOr(err error, notFoundMsg, internalMsg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NotFound(notFoundMsg)
	}
	return apperror.Internal(internalMsg, err)
}
