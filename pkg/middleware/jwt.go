package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/musobarlab/rumpi/pkg/shared"
)

// ValidateJWT function
func (m *Middleware) ValidateJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authorization := c.Request().Header.Get(echo.HeaderAuthorization)
			if authorization == "" {
				return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid authorization").JSON(c.Response())
			}

			authValues := strings.Split(authorization, " ")
			authType := strings.ToLower(authValues[0])
			if authType != "bearer" || len(authValues) != 2 {
				return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid authorization").JSON(c.Response())
			}

			tokenString := authValues[1]
			resp := m.jwtService.Validate(c.Request().Context(), tokenString)
			if resp.Error != nil {
				return shared.NewHTTPResponse(http.StatusUnauthorized, resp.Error.Error()).JSON(c.Response())
			}

			c.Set("jwtClaim", resp.Data)
			return next(c)
		}
	}
}
