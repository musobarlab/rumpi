package delivery

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	"github.com/musobarlab/rumpi/internal/modules/user/domain"
	"github.com/musobarlab/rumpi/internal/modules/user/usecase"
	"github.com/musobarlab/rumpi/pkg/chathub"
	"github.com/musobarlab/rumpi/pkg/jwt"
	"github.com/musobarlab/rumpi/pkg/middleware"
	"github.com/musobarlab/rumpi/pkg/shared"
)

// EchoDelivery represent http handler handled by Echo
type EchoDelivery struct {
	userUsecase usecase.UserUsecase
	middleware  *middleware.Middleware
	chatManager *chathub.Manager
}

func NewEchoDelivery(userUsecase usecase.UserUsecase, middleware *middleware.Middleware, chatManager *chathub.Manager) *EchoDelivery {
	return &EchoDelivery{
		userUsecase: userUsecase,
		middleware:  middleware,
		chatManager: chatManager,
	}
}

// Mount function will mount each EchiDelivery into route group
func (h *EchoDelivery) Mount(root *echo.Group) {
	root.POST("/login", h.login, h.middleware.BasicAuth())
	root.POST("/register", h.register, h.middleware.BasicAuth())
	root.GET("/chat", h.chat)
	root.GET("/profile", h.getProfile, h.middleware.ValidateJWT())
}

func (h *EchoDelivery) login(c echo.Context) error {

	var loginRequest domain.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return shared.NewHTTPResponse(http.StatusInternalServerError, "invalid payload").JSON(c.Response())
	}

	ipAddress := c.RealIP()
	ctx := shared.SetToContext(c.Request().Context(), shared.ContextKey("ipAddress"), ipAddress)
	loginResult := h.userUsecase.Login(ctx, &loginRequest)
	if loginResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid username or password").JSON(c.Response())
	}

	loginResponse := loginResult.Data.(*domain.LoginResponse)

	return shared.NewHTTPResponse(http.StatusOK, "login success", loginResponse).JSON(c.Response())
}

func (h *EchoDelivery) register(c echo.Context) error {

	var userRequest domain.User
	if err := c.Bind(&userRequest); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid user payload").JSON(c.Response())
	}

	ipAddress := c.RealIP()
	ctx := shared.SetToContext(c.Request().Context(), shared.ContextKey("ipAddress"), ipAddress)
	registerResult := h.userUsecase.Register(ctx, &userRequest)
	if registerResult.Error != nil {
		if registerResult.Error == shared.ErrUserAlreadyExist {
			return shared.NewHTTPResponse(http.StatusBadRequest, registerResult.Error.Error()).JSON(c.Response())
		}

		return shared.NewHTTPResponse(http.StatusBadRequest, "error processing your request").JSON(c.Response())
	}

	userResponse := registerResult.Data.(*domain.User)

	return shared.NewHTTPResponse(http.StatusCreated, "register success", userResponse).JSON(c.Response())
}

func (h *EchoDelivery) getProfile(c echo.Context) error {

	ipAddress := c.RealIP()
	ctx := shared.SetToContext(c.Request().Context(), shared.ContextKey("ipAddress"), ipAddress)

	jwtClaimCtx := c.Get("jwtClaim")
	jwtClaim := jwtClaimCtx.(*jwt.Claim)

	userID, err := primitive.ObjectIDFromHex(jwtClaim.User.ID)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid user id").JSON(c.Response())
	}

	profileResult := h.userUsecase.GetProfile(ctx, &domain.User{ID: userID})
	if profileResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, "error processing your request").JSON(c.Response())
	}

	userResponse := profileResult.Data.(*domain.User)

	return shared.NewHTTPResponse(http.StatusOK, "get profile success", userResponse).JSON(c.Response())
}

func (h *EchoDelivery) chat(c echo.Context) error {
	sock, err := h.chatManager.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	id := c.Request().Header.Get("Sec-Websocket-Key")

	var client chathub.Client
	client.ID = id
	client.Conn = sock
	client.MsgChan = make(chan []byte)
	client.Room = make(map[string]bool)
	client.Manager = h.chatManager

	h.chatManager.Register <- &client

	// Consume message
	go client.Consume()

	// Publish message
	go client.Publish()

	return nil
}
