package main

import (
	"fmt"
	"log"
	"os"
	"saudemais-api/internal/database"
	"saudemais-api/internal/router"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	e := echo.New()

	router.SetupRoutes(e, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Servidor rodando na porta", port)
	e.Logger.Fatal(e.Start(":" + port))
}
