package inventory

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventoryRepo "github.com/ahmadammarm/amami/backend/internal/repository/inventory"
	"github.com/google/uuid"
)

type InventoryService interface {
	CreateAsset(req inventoryDto.CreateAssetRequest) (*inventoryDto.AssetResponse, error)
	GetAllAssets(page, limit int) ([]domain.Asset, int64, error)
	UpdateAsset(id string, req inventoryDto.CreateAssetRequest) (*inventoryDto.AssetResponse, error)
	DeleteAsset(id string) error

	CreateLoan(assetID uuid.UUID, req inventoryDto.CreateAssetLoanRequest) error
	ReturnLoan(loanID uuid.UUID, req inventoryDto.ReturnAssetLoanRequest) error
}

type inventoryService struct {
	repo inventoryRepo.InventoryRepository
}

func NewInventoryService(repo inventoryRepo.InventoryRepository) InventoryService {
	return &inventoryService{repo}
}

func (s *inventoryService) CreateAsset(req inventoryDto.CreateAssetRequest) (*inventoryDto.AssetResponse, error) {
	asset := &domain.Asset{
		Name:          req.Name,
		SKU:           req.SKU,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		CurrentStatus: req.CurrentStatus,
		Location:      req.Location,
	}

	if err := s.repo.CreateAsset(asset); err != nil {
		return nil, err
	}

	return &inventoryDto.AssetResponse{
		ID:            asset.ID,
		Name:          asset.Name,
		SKU:           asset.SKU,
		PurchaseDate:  asset.PurchaseDate,
		PurchasePrice: asset.PurchasePrice,
		CurrentStatus: asset.CurrentStatus,
		Location:      asset.Location,
	}, nil
}

func (s *inventoryService) GetAllAssets(page, limit int) ([]domain.Asset, int64, error) {
	return s.repo.GetAllAssets(page, limit)
}

func (s *inventoryService) UpdateAsset(id string, req inventoryDto.CreateAssetRequest) (*inventoryDto.AssetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}
	asset, err := s.repo.GetAssetByID(uid)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	asset.Name = req.Name
	asset.SKU = req.SKU
	asset.PurchaseDate = req.PurchaseDate
	asset.PurchasePrice = req.PurchasePrice
	asset.CurrentStatus = req.CurrentStatus
	asset.Location = req.Location

	if err := s.repo.UpdateAsset(asset); err != nil {
		return nil, err
	}

	return &inventoryDto.AssetResponse{
		ID:            asset.ID,
		Name:          asset.Name,
		SKU:           asset.SKU,
		PurchaseDate:  asset.PurchaseDate,
		PurchasePrice: asset.PurchasePrice,
		CurrentStatus: asset.CurrentStatus,
		Location:      asset.Location,
	}, nil
}

func (s *inventoryService) DeleteAsset(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteAsset(uid)
}

func (s *inventoryService) CreateLoan(assetID uuid.UUID, req inventoryDto.CreateAssetLoanRequest) error {
	asset, err := s.repo.GetAssetByID(assetID)
	if err != nil {
		return err
	}
	if asset.CurrentStatus != "GOOD" {
		return errors.New("only GOOD assets can be loaned")
	}

	loan := &domain.AssetLoan{
		AssetID:        assetID,
		JamaahID:       req.JamaahID,
		LoanDate:       req.LoanDate,
		DueDate:        req.DueDate,
		ConditionNotes: req.ConditionNotes,
	}
	
	if err := s.repo.CreateLoan(loan); err != nil {
		return err
	}
	
	asset.CurrentStatus = "LOANED"
	return s.repo.UpdateAsset(asset)
}

func (s *inventoryService) ReturnLoan(loanID uuid.UUID, req inventoryDto.ReturnAssetLoanRequest) error {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		return err
	}
	loan.ReturnDate = &req.ReturnDate
	if req.ConditionNotes != "" {
		loan.ConditionNotes = loan.ConditionNotes + " | Return: " + req.ConditionNotes
	}
	
	if err := s.repo.UpdateLoan(loan); err != nil {
		return err
	}
	
	asset, err := s.repo.GetAssetByID(loan.AssetID)
	if err == nil {
		asset.CurrentStatus = "GOOD"
		s.repo.UpdateAsset(asset)
	}
	return nil
}
