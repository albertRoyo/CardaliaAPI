/*
File		: main.go
Description	: Main file of the project. The main connects to the database and defines the router and all of its endpoints and CORS
The CardaliaAPI diagram from the memory of the project (Figure 8) is very usefull to understand all the backend of the aplication.
*/

package main

import (
	"CardaliaAPI/middlewares"
	"CardaliaAPI/models"
	"CardaliaAPI/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Methods
	router.POST("/register", routes.Register)
	router.POST("/login", routes.Login)

	router.GET("/cards/:autocomplete", routes.GetCardsByName)
	router.GET("/cards/versions/:cardname", routes.GetCardVersions)

	router.GET("/user/collection/:username", routes.GetUserCollectionByName)

	protected := router.Group("/")
	protected.Use(middlewares.JwtAuthMiddleware())

	// Private methods
	protected.PUT("/user/password", routes.ChangeUserPassword)

	protected.POST("/user/collection", routes.SaveCollection)
	protected.GET("/user/collection", routes.GetCollection)

	protected.GET("/users/collections/:card_id", routes.GetAllUserCollectionsByCardId)

	protected.POST("/user/trade", routes.NewTrade)
	protected.PUT("/user/trade", routes.ModifyTrade)
	protected.DELETE("/user/trade/:username", routes.DeleteTrade)

	protected.GET("/user/trades", routes.GetTrades)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" || port == "" {
		log.Fatal("$PORT and $HOST must be set")
	}
	router.Run(host + ":" + port)

}
