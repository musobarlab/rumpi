package server

import (
	"fmt"
	"net/http"

	"github.com/musobarlab/rumpi/pkg/shared"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/musobarlab/rumpi/config"
	uDelivery "github.com/musobarlab/rumpi/internal/modules/user/delivery"
)

// HTTPServer struct represent http server
type HTTPServer struct {
	UserEchoDelivery *uDelivery.EchoDelivery
}

func (s *HTTPServer) Run() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/", func(c echo.Context) error {
		return shared.NewHTTPResponse(http.StatusOK, "server up and running").JSON(c.Response())
	})

	userGroup := e.Group("/users")
	s.UserEchoDelivery.Mount(userGroup)

	return e.Start(fmt.Sprintf(":%d", config.Config.HTTPPort))
}
