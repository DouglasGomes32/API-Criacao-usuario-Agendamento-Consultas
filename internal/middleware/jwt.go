package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AutenticarJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Token ausente ou inválido",
					"code":    401,
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "segredo_padrao"
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Método de assinatura inválido")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Token inválido ou expirado",
					"code":    401,
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["sub"] == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Token inválido",
					"code":    401,
				})
			}

			c.Set("userEmail", claims["sub"])

			return next(c)
		}
	}
}
