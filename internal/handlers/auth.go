package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"saudemais-api/internal/models"
	"saudemais-api/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req types.RegisterRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Dados inválidos",
				Code:    400,
			})
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Nome, email e senha são obrigatórios",
				Code:    400,
			})
		}

		err := models.CriarPaciente(db, req.Name, req.Email, req.Password)
		if err != nil {
			if err.Error() == "Paciente já cadastrado" {
				return c.JSON(http.StatusConflict, types.APIResponse{
					Message: err.Error(),
					Code:    409,
				})
			}
			log.Println("Erro ao registrar paciente:", err)
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao registrar paciente",
				Code:    500,
			})
		}

		return c.JSON(http.StatusCreated, types.APIResponse{
			Message: "Paciente registrado com sucesso",
			Code:    201,
		})
	}
}

func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req types.LoginRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Dados inválidos",
				Code:    400,
			})
		}

		if req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Email e senha são obrigatórios",
				Code:    400,
			})
		}

		paciente, err := models.BuscarPacientePorEmail(db, req.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Paciente não encontrado",
				Code:    401,
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(paciente.PasswordHash), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Senha incorreta",
				Code:    401,
			})
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "segredo_padrao"
		}

		claims := &jwt.RegisteredClaims{
			Subject:   paciente.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao gerar token",
				Code:    500,
			})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": t})
	}
}
