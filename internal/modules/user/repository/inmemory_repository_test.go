package repository

import (
	"context"
	"testing"

	"crypto/sha1"

	p "github.com/wuriyanto48/go-pbkdf2"

	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInMemoryRpository(t *testing.T) {

	db := make(map[primitive.ObjectID]*domain.User)
	idDummy, _ := primitive.ObjectIDFromHex("5e9aa00562d663ea2509c0ee")
	db[idDummy] = &domain.User{
		FullName: "Alex",
		Email:    "alex@yahoo.co.id",
		Photo:    "https://storge.rumpi.com/profile.png",
	}

	userRepository := NewInmemoryRepository(db)

	t.Run("test save user", func(t *testing.T) {
		pass := p.NewPassword(sha1.New, 8, 32, 1000)
		hashed := pass.HashPassword("12345")

		user := &domain.User{
			FullName: "Wuriyanto",
			Email:    "wuriyanto@yahoo.co.id",
			Password: hashed.CipherText,
			Salt:     hashed.Salt,
			Photo:    "https://storge.rumpi.com/profile.png",
		}

		result := userRepository.Save(context.Background(), user)

		if result.Data == nil {
			t.Error("should success save user")
		}

	})

	t.Run("test find by id", func(t *testing.T) {
		result := userRepository.Find(context.Background(), &domain.User{ID: idDummy})
		if result.Error != nil {
			t.Error("should success find user by id")
		}

		if result.Data == nil {
			t.Error("user should not nil")
		}

		expected := "alex@yahoo.co.id"

		user := result.Data.(*domain.User)
		if user.Email != expected {
			t.Error("user's email should equal alex@yahoo.co.id")
		}
	})

	t.Run("test find by email", func(t *testing.T) {
		result := userRepository.Find(context.Background(), &domain.User{Email: "alex@yahoo.co.id"})
		if result.Error != nil {
			t.Error("should success find user by email")
		}

		if result.Data == nil {
			t.Error("user should not nil")
		}

		expected := "alex@yahoo.co.id"

		user := result.Data.(*domain.User)
		if user.Email != expected {
			t.Error("user's email should equal alex@yahoo.co.id")
		}
	})
}
