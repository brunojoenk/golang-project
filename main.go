package main

import (
	"fmt"
	"github/brunojoenk/golang-test/config"
	"github/brunojoenk/golang-test/database"
	"github/brunojoenk/golang-test/handlers"
	"github/brunojoenk/golang-test/models"

	_ "github/brunojoenk/golang-test/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Swagger API
// @version 1.0
// @description This is a sample server to manager books.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := config.New()

	database, err := database.NewPsqlDB(cfg)
	if err != nil {
		e.Logger.Fatal("Error on open database connectiont: ", err.Error())
	}

	// Create/update tables on database
	err = database.AutoMigrate(&models.Author{}, &models.Book{})
	if err != nil {
		e.Logger.Fatal("Error on execute migrate: ", err.Error())
	}

	h := handlers.New(database)
	h.HandleControllers(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ServerPort)))

}
