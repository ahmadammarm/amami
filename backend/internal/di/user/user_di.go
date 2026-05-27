package user

import (
	userHandler "github.com/ahmadammarm/amami/backend/internal/handler/user"
	authRepo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	userRepo "github.com/ahmadammarm/amami/backend/internal/repository/user"
	userSvc "github.com/ahmadammarm/amami/backend/internal/service/user"
	"gorm.io/gorm"
)

type Module struct {
	Handler userHandler.UserMgmtHandler
}

func NewModule(db *gorm.DB, roleRepo authRepo.RoleRepository) *Module {
	repo := userRepo.NewUserMgmtRepository(db)
	svc := userSvc.NewUserMgmtService(repo, roleRepo)
	h := userHandler.NewUserMgmtHandler(svc)
	
	return &Module{Handler: h}
}
