package usecase

import (
	"context"

	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"github.com/musobarlab/rumpi/pkg/shared"
)

// UserUsecase usecase
type UserUsecase interface {
	Register(ctx context.Context, user *domain.User) shared.Result
	Login(ctx context.Context, userLogin *domain.LoginRequest) shared.Result
	IsUserExist(ctx context.Context, email string) shared.Result
	GetProfile(ctx context.Context, user *domain.User) shared.Result
}
