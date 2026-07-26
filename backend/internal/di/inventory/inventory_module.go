package inventory

import (
	inventoryHandler "github.com/ahmadammarm/amami/backend/internal/handler/inventory"
	inventoryRepo "github.com/ahmadammarm/amami/backend/internal/repository/inventory"
	inventorySvc "github.com/ahmadammarm/amami/backend/internal/service/inventory"
	"gorm.io/gorm"
)

type Module struct {
	Handler *inventoryHandler.InventoryHandler
}

func NewModule(db *gorm.DB) *Module {
	repo := inventoryRepo.NewInventoryRepository(db)
	svc := inventorySvc.NewInventoryService(repo)
	handler := inventoryHandler.NewInventoryHandler(svc)

	return &Module{
		Handler: handler,
	}
}
