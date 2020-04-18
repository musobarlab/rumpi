package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/musobarlab/rumpi/config"
	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"github.com/musobarlab/rumpi/internal/modules/user/repository"
	"github.com/musobarlab/rumpi/pkg/jwt"
	"github.com/musobarlab/rumpi/pkg/shared"
	p "github.com/wuriyanto48/go-pbkdf2"
)

// UserUsecaseImpl implementation of UserUsecase
type UserUsecaseImpl struct {
	userRepository repository.UserRepository
	passwordHasher *p.Password
	jwtService     jwt.JwtService
}

// NewUserUsecaseImpl function will initialize UserUsecaseImpl
func NewUserUsecaseImpl(userRepository repository.UserRepository, passwordHasher *p.Password, jwtService jwt.JwtService) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		jwtService:     jwtService,
	}
}

// Register function will register new user
func (u *UserUsecaseImpl) Register(ctx context.Context, user *domain.User) shared.Result {

	resultExist := u.IsUserExist(ctx, user.Email)
	if resultExist.Error != nil {
		return shared.Result{Error: resultExist.Error}
	}

	exist := resultExist.Data.(bool)
	if exist {
		return shared.Result{Error: shared.ErrUserAlreadyExist}
	}

	hashed := u.passwordHasher.HashPassword(user.Password)
	user.Password = hashed.CipherText
	user.Salt = hashed.Salt

	result := u.userRepository.Save(ctx, user)
	if resultExist.Error != nil {
		return shared.Result{Error: result.Error}
	}

	return shared.Result{Data: result.Data}
}

// Login function will help user logging in into system
func (u *UserUsecaseImpl) Login(ctx context.Context, userLogin *domain.LoginRequest) shared.Result {
	user := &domain.User{Email: userLogin.Username}
	result := u.userRepository.Find(ctx, user)
	if result.Error != nil {
		return shared.Result{Error: result.Error}
	}

	user = result.Data.(*domain.User)
	if !user.IsValidPassword(userLogin.Password) {
		return shared.Result{Error: shared.ErrInvalidUsernameOrPassword}
	}

	var jwtClaim jwt.Claim
	jwtClaim.Alg = jwt.RS256
	jwtClaim.Subject = user.ID.Hex()
	jwtClaim.User.ID = user.ID.Hex()
	jwtClaim.User.Email = user.Email
	jwtClaim.User.FullName = user.FullName

	// TODO
	var accessTokenExpired time.Duration
	if userLogin.AccessTokenExpired > 0 {
		var err error
		accessTokenExpiredStr := fmt.Sprintf("%dm", userLogin.AccessTokenExpired)
		accessTokenExpired, err = time.ParseDuration(accessTokenExpiredStr)
		if err != nil {
			accessTokenExpired = config.Config.AccessTokenExpired
		}
	} else {
		accessTokenExpired = config.Config.AccessTokenExpired
	}

	jwtResult := u.jwtService.Generate(ctx, &jwtClaim, accessTokenExpired)
	if jwtResult.Error != nil {
		return shared.Result{Error: jwtResult.Error}
	}

	accessToken, ok := jwtResult.Data.(string)
	if !ok {
		return shared.Result{Error: shared.ErrTokenFormat}
	}

	return shared.Result{Data: &domain.LoginResponse{
		Email:              user.Email,
		AccessToken:        accessToken,
		AccessTokenExpired: int64(accessTokenExpired.Minutes()),
	}}
}

// IsUserExist function will check if user already exist
func (u *UserUsecaseImpl) IsUserExist(ctx context.Context, email string) shared.Result {
	result := u.userRepository.Find(ctx, &domain.User{Email: email})
	if result.Error != nil {
		if _, ok := result.Error.(*shared.ErrorUserNotFound); ok {
			return shared.Result{Data: false}
		}
		return shared.Result{Error: result.Error}
	}

	return shared.Result{Data: true}
}

// GetProfile function will get profile base on its identifier
func (u *UserUsecaseImpl) GetProfile(ctx context.Context, user *domain.User) shared.Result {
	result := u.userRepository.Find(ctx, user)
	if result.Error != nil {
		return shared.Result{Error: result.Error}
	}

	return shared.Result{Data: result.Data}
}
