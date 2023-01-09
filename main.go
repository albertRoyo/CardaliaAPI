package main

import (
	"CardaliaAPI/models"
	"CardaliaAPI/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type") //"Access-Control-Allow-Headers, Origin, Accept, Authorization, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	models.ConnectDataBase()
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//public := router.Group("/")

	//router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/register", routes.Register)
	router.POST("/login", routes.Login)

	router.GET("/card/:cardname", routes.GetCardByName)
	router.GET("/:autocomplete", routes.GetCardsByName)

	router.GET("/cards/versions/:cardname", routes.GetCardVersions)
	router.GET("/versions/:cardname", routes.GetVersionNames)

	router.GET("/cards/:set/:number", routes.GetCardVersion)

	router.GET("/user/:username", routes.GetUserCollectionByName)
	router.GET("/users/:card_id", routes.GetAllUserCollectionsByCardId)

	//protected := router.Group("/admin")
	//protected.Use(cors.Default())
	//protected.Use(middlewares.JwtAuthMiddleware())
	router.POST("/cards", routes.SaveCollection)
	router.GET("/cards", routes.GetCollection)

	router.POST("/trade", routes.NewTrade)
	router.PUT("/trade", routes.ModifyTrade)
	router.DELETE("/trade/:username", routes.DeleteTrade)

	router.GET("/trades", routes.GetTrades)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router.Run(":" + port)

}
