package auth

import (
	handler "github.com/ahmadammarm/amami/backend/internal/handler/auth"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/auth"
	svc "github.com/ahmadammarm/amami/backend/internal/service/auth"
	"github.com/google/wire"
)

// ProviderSet is the provider set for the Auth module
var ProviderSet = wire.NewSet(
	repo.NewUserRepository,
	repo.NewRoleRepository,
	svc.NewAuthService,
	handler.NewAuthHandler,
)
