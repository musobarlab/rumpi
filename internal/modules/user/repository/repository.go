package repository

import (
	"context"

	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"github.com/musobarlab/rumpi/pkg/shared"
)

// UserRepository repository
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) shared.Result
	Find(ctx context.Context, where *domain.User) shared.Result
	Count(ctx context.Context, where *domain.User) shared.Result
}
