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
		return nil, toAppError(err, "failed to create asset")
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
	asset, err := s.repo.UpdateAsset(id, editAsset(req))
	if err != nil {
		return nil, toAppError(err, "failed to update asset")
	}

	res := inventoryDto.ToAssetResponse(asset)
	return &res, nil
}

func (s *inventoryService) DeleteAsset(id uuid.UUID) error {
	if err := s.repo.DeleteAsset(id, ensureNotOnLoan); err != nil {
		return toAppError(err, "failed to delete asset")
	}
	return nil
}

func (s *inventoryService) CreateLoan(assetID uuid.UUID, req inventoryDto.CreateAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error) {
	// Input-only checks run before any lock is taken.
	if req.DueDate != nil && req.DueDate.Before(req.LoanDate) {
		return nil, apperror.Validation("due date must not be before the loan date")
	}

	loan, err := s.repo.CreateLoan(assetID, newLoan(req))
	if err != nil {
		return nil, toAppError(err, "failed to create loan")
	}

	res := inventoryDto.ToLoanResponse(loan)
	return &res, nil
}

func (s *inventoryService) ReturnLoan(loanID uuid.UUID, req inventoryDto.ReturnAssetLoanRequest) (*inventoryDto.AssetLoanResponse, error) {
	loan, err := s.repo.ReturnLoan(loanID, recordReturn(req))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("loan not found")
		}
		return nil, toAppError(err, "failed to return loan")
	}

	res := inventoryDto.ToLoanResponse(loan)
	return &res, nil
}

// Business rules. Each runs inside the repository transaction against the
// locked row, so the state it checks cannot change before the write.

// editAsset applies the request to an asset. LOANED is managed by the loan
// workflow, so an edit can neither set it nor clear it.
func editAsset(req inventoryDto.AssetRequest) inventoryRepo.AssetEdit {
	return func(asset *domain.Asset) error {
		if (asset.CurrentStatus == domain.AssetLoaned) != (req.CurrentStatus == domain.AssetLoaned) {
			return apperror.Conflict("loan status is managed through loans; create or return a loan instead")
		}
		applyAssetRequest(asset, req)
		return nil
	}
}

// ensureNotOnLoan rejects deleting an asset that is currently lent out.
func ensureNotOnLoan(asset *domain.Asset) error {
	if asset.CurrentStatus == domain.AssetLoaned {
		return apperror.Conflict("asset is on loan; return it before deleting")
	}
	return nil
}

// newLoan lends an asset that is in good condition.
func newLoan(req inventoryDto.CreateAssetLoanRequest) inventoryRepo.LoanBuilder {
	return func(asset *domain.Asset) (*domain.AssetLoan, error) {
		if asset.CurrentStatus != domain.AssetGood {
			return nil, apperror.Conflict("only assets in GOOD condition can be loaned")
		}
		return &domain.AssetLoan{
			AssetID:        asset.ID,
			JamaahID:       req.JamaahID,
			LoanDate:       req.LoanDate,
			DueDate:        req.DueDate,
			ConditionNotes: req.ConditionNotes,
		}, nil
	}
}

// recordReturn closes an open loan.
func recordReturn(req inventoryDto.ReturnAssetLoanRequest) inventoryRepo.LoanReturn {
	return func(loan *domain.AssetLoan) error {
		if loan.ReturnDate != nil {
			return apperror.Conflict("loan has already been returned")
		}
		if req.ReturnDate.Before(loan.LoanDate) {
			return apperror.Validation("return date must not be before the loan date")
		}
		loan.ReturnDate = &req.ReturnDate
		if req.ConditionNotes != "" {
			loan.ConditionNotes = loan.ConditionNotes + " | Return: " + req.ConditionNotes
		}
		return nil
	}
}

func applyAssetRequest(asset *domain.Asset, req inventoryDto.AssetRequest) {
	asset.Name = req.Name
	asset.SKU = req.SKU
	asset.PurchaseDate = req.PurchaseDate
	asset.PurchasePrice = req.PurchasePrice
	asset.CurrentStatus = req.CurrentStatus
	asset.Location = req.Location
}

func toAppError(err error, internalMsg string) error {
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return apperror.Conflict("an asset with this SKU already exists")
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return apperror.Validation("borrower not found")
	default:
		return apperror.FromStore(err, "asset not found", internalMsg)
	}
}
