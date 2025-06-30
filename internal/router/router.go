package router

import (
	"database/sql"

	"saudemais-api/internal/handlers"
	"saudemais-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *sql.DB) {
	// Rotas públicas
	e.POST("/register", handlers.Register(db))
	e.POST("/login", handlers.Login(db))

	// Rotas protegidas
	jwtMid := middleware.AutenticarJWT()

	appointments := e.Group("/appointments")
	appointments.Use(jwtMid)
	appointments.POST("", handlers.AgendarConsulta(db))
	appointments.GET("", handlers.ListarConsulta(db))
	appointments.DELETE("/:id", handlers.CancelarConsulta(db))
}
