package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/musobarlab/rumpi/pkg/shared"
)

// BasicAuth function basic auth
func (m *Middleware) BasicAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizations := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(authorizations) != 2 {
				return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid app id and app key").JSON(c.Response())
			}

			authType, val := authorizations[0], authorizations[1]
			if authType != "Basic" {
				return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid app id and app key").JSON(c.Response())
			}

			isValid := func() bool {
				data, err := base64.StdEncoding.DecodeString(val)
				if err != nil {
					return false
				}

				decoded := strings.Split(string(data), ":")
				if len(decoded) < 2 {
					return false
				}
				username, password := decoded[0], decoded[1]

				if username != m.username || password != m.password {
					return false
				}

				return true
			}

			if !isValid() {
				return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid app id and app key").JSON(c.Response())
			}

			return next(c)
		}
	}
}
