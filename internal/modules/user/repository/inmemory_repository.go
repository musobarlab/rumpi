package repository

import (
	"context"
	"errors"
	"time"

	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"github.com/musobarlab/rumpi/pkg/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InmemoryRepository in memory implementation of UserRepository
type InmemoryRepository struct {
	db map[primitive.ObjectID]*domain.User
}

// NewInmemoryRepository function will initialize InmemoryRepository
func NewInmemoryRepository(db map[primitive.ObjectID]*domain.User) *InmemoryRepository {
	return &InmemoryRepository{db}
}

// Save function, function will store user data into InmemoryRepository
func (r *InmemoryRepository) Save(ctx context.Context, user *domain.User) shared.Result {
	if user == nil {
		return shared.Result{Error: errors.New("user cannot be empty")}
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	user.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	r.db[user.ID] = user

	return shared.Result{Data: user}
}

// Find function, function will find user data from InmemoryRepository with where filter
func (r *InmemoryRepository) Find(ctx context.Context, where *domain.User) shared.Result {
	if !where.ID.IsZero() {
		user, ok := r.db[where.ID]
		if !ok {
			return shared.Result{Error: &shared.ErrorUserNotFound{Message: "user not found"}}
		}

		return shared.Result{Data: user}

	}

	if where.Email != "" {
		for _, user := range r.db {
			if where.Email == user.Email {
				return shared.Result{Data: user}
			}
		}
	}
	return shared.Result{Error: &shared.ErrorUserNotFound{Message: "user not found"}}
}

// Count function, function will count user data from InmemoryRepository with what filter
func (r *InmemoryRepository) Count(ctx context.Context, where *domain.User) shared.Result {
	return shared.Result{}
}
