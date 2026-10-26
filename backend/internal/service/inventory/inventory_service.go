package inventory

import (
	"errors"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	inventoryDto "github.com/ahmadammarm/amami/backend/internal/dto/inventory"
	inventoryRepo "github.com/ahmadammarm/amami/backend/internal/repository/inventory"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryService interface {
	CreateAsset(req inventoryDto.AssetRequest) (*inventoryDto.AssetResponse, error)
	GetAllAssets(page, limit int) (utils.Page[inventoryDto.AssetResponse], error)
	UpdateAsset(id uuid.UUID, req inventoryDto.AssetRequest) (*inventoryDto.AssetResponse, error)
	DeleteAsset(id uuid.UUID) error

	CreateLoan(assetID uuid.UUID, req inventoryDto.CreateAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error)
	ReturnLoan(loanID uuid.UUID, req inventoryDto.ReturnAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error)
}

type inventoryService struct {
	repo inventoryRepo.InventoryRepository
}

func NewInventoryService(repo inventoryRepo.InventoryRepository) InventoryService {
	return &inventoryService{repo}
}

func (s *inventoryService) CreateAsset(req inventoryDto.AssetRequest) (*inventoryDto.AssetResponse, error) {
	if req.CurrentStatus == domain.AssetLoaned {
		return nil, apperror.Validation("an asset becomes LOANED by creating a loan")
	}

	asset := &domain.Asset{}
	applyAssetRequest(asset, req)

	if err := s.repo.CreateAsset(asset); err != nil {
		return nil, mapAssetWriteError(err, "failed to create asset")
	}

	res := inventoryDto.ToAssetResponse(asset)
	return &res, nil
}

func (s *inventoryService) GetAllAssets(page, limit int) (utils.Page[inventoryDto.AssetResponse], error) {
	assets, total, err := s.repo.GetAllAssets(page, limit)
	if err != nil {
		return utils.Page[inventoryDto.AssetResponse]{}, apperror.Internal("failed to fetch assets", err)
	}

	res := make([]inventoryDto.AssetResponse, 0, len(assets))
	for i := range assets {
		res = append(res, inventoryDto.ToAssetResponse(&assets[i]))
	}
	return utils.NewPage(res, total, page, limit), nil
}

func (s *inventoryService) UpdateAsset(id uuid.UUID, req inventoryDto.AssetRequest) (*inventoryDto.AssetResponse, error) {
	asset, err := s.findAsset(id)
	if err != nil {
		return nil, err
	}

	// LOANED is managed by the loan workflow, so it can neither be set nor
	// cleared by editing the asset directly.
	if (asset.CurrentStatus == domain.AssetLoaned) != (req.CurrentStatus == domain.AssetLoaned) {
		return nil, apperror.Conflict("loan status is managed through loans; create or return a loan instead")
	}

	applyAssetRequest(asset, req)
	if err := s.repo.UpdateAsset(asset); err != nil {
		return nil, mapAssetWriteError(err, "failed to update asset")
	}

	res := inventoryDto.ToAssetResponse(asset)
	return &res, nil
}

func (s *inventoryService) DeleteAsset(id uuid.UUID) error {
	asset, err := s.findAsset(id)
	if err != nil {
		return err
	}
	if asset.CurrentStatus == domain.AssetLoaned {
		return apperror.Conflict("asset is on loan; return it before deleting")
	}

	if err := s.repo.DeleteAsset(id); err != nil {
		return apperror.Internal("failed to delete asset", err)
	}
	return nil
}

func (s *inventoryService) CreateLoan(assetID uuid.UUID, req inventoryDto.CreateAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error) {
	asset, err := s.findAsset(assetID)
	if err != nil {
		return nil, err
	}
	if asset.CurrentStatus != domain.AssetGood {
		return nil, apperror.Conflict("only assets in GOOD condition can be loaned")
	}
	if req.DueDate != nil && req.DueDate.Before(req.LoanDate) {
		return nil, apperror.Validation("due date must not be before the loan date")
	}

	loan := &domain.AssetLoan{
		AssetID:        assetID,
		JamaahID:       req.JamaahID,
		LoanDate:       req.LoanDate,
		DueDate:        req.DueDate,
		ConditionNotes: req.ConditionNotes,
	}

	if err := s.repo.CreateLoan(loan); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, apperror.Validation("borrower not found")
		}
		return nil, apperror.Internal("failed to create loan", err)
	}

	res := inventoryDto.ToLoanResponse(loan)
	return &res, nil
}

func (s *inventoryService) ReturnLoan(loanID uuid.UUID, req inventoryDto.ReturnAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error) {
	loan, err := s.repo.GetLoanByID(loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("loan not found")
		}
		return nil, apperror.Internal("failed to fetch loan", err)
	}
	if loan.ReturnDate != nil {
		return nil, apperror.Conflict("loan has already been returned")
	}
	if req.ReturnDate.Before(loan.LoanDate) {
		return nil, apperror.Validation("return date must not be before the loan date")
	}

	loan.ReturnDate = &req.ReturnDate
	if req.ConditionNotes != "" {
		loan.ConditionNotes = loan.ConditionNotes + " | Return: " + req.ConditionNotes
	}

	if err := s.repo.ReturnLoan(loan); err != nil {
		return nil, apperror.Internal("failed to return loan", err)
	}

	res := inventoryDto.ToLoanResponse(loan)
	return &res, nil
}

func (s *inventoryService) findAsset(id uuid.UUID) (*domain.Asset, error) {
	asset, err := s.repo.GetAssetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("asset not found")
		}
		return nil, apperror.Internal("failed to fetch asset", err)
	}
	return asset, nil
}

func applyAssetRequest(asset *domain.Asset, req inventoryDto.AssetRequest) {
	asset.Name = req.Name
	asset.SKU = req.SKU
	asset.PurchaseDate = req.PurchaseDate
	asset.PurchasePrice = req.PurchasePrice
	asset.CurrentStatus = req.CurrentStatus
	asset.Location = req.Location
}

func mapAssetWriteError(err error, msg string) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.Conflict("an asset with this SKU already exists")
	}
	return apperror.Internal(msg, err)
}
