package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"

	"github.com/musobarlab/rumpi/config"
	userDelivery "github.com/musobarlab/rumpi/internal/modules/user/delivery"
	userDomain "github.com/musobarlab/rumpi/internal/modules/user/domain"
	userRepo "github.com/musobarlab/rumpi/internal/modules/user/repository"
	userUc "github.com/musobarlab/rumpi/internal/modules/user/usecase"
	"github.com/musobarlab/rumpi/internal/server"
	"github.com/musobarlab/rumpi/pkg/chathub"
	"github.com/musobarlab/rumpi/pkg/jwt"
	"github.com/musobarlab/rumpi/pkg/middleware"
	p "github.com/wuriyanto48/go-pbkdf2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("starting application...")

	passwordHasher := p.NewPassword(sha1.New, 8, 32, int(config.Config.PasswordHashIteration))

	jwtService := jwt.NewJWT(config.Config.PublicKey, config.Config.PrivateKey)

	chatManager := chathub.NewManager(config.Config.WebsocketKey, jwtService)

	// user module
	inMemoryDB := make(map[primitive.ObjectID]*userDomain.User)
	userRepository := userRepo.NewInmemoryRepository(inMemoryDB)

	mw := middleware.NewMiddleware()

	userUsecase := userUc.NewUserUsecaseImpl(userRepository, passwordHasher, jwtService)
	userEchoDelivery := userDelivery.NewEchoDelivery(userUsecase, mw, chatManager)

	httpServer := &server.HTTPServer{
		UserEchoDelivery: userEchoDelivery,
	}

	go chatManager.Handle()

	log.Fatal(httpServer.Run())
}
