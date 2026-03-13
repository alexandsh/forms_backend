package main

// @title           SL Forms API
// @version         1.0
// @description     API form-сервиса
// @BasePath        /

import (
	"context"
	"log"

	_ "forms/docs"
	"forms/internal/db"
	"forms/internal/handler"
	"forms/internal/middleware"
	"forms/internal/repository"
	"forms/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	godotenv.Load()

	router := gin.Default()

	ctx := context.Background()

	pool, err := db.ConnectToDB(ctx)
	if err != nil {
		log.Fatalf("connetction failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("connetction failed: %v", err)
	}

	log.Println("connection done")

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	formRepo := repository.NewFormRepository(pool)
	formService := service.NewFormService(formRepo)
	formHandler := handler.NewFormHandler(formService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware())
	{
		api.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID})
		})
		api.POST("/form", formHandler.CreateForm)
		api.GET("/form/:id", formHandler.GetForm)
		api.GET("/forms", formHandler.GetForms)
		api.PATCH("/form/:id", formHandler.UpdateForm)
		api.DELETE("/form/:id", formHandler.DeleteForm)
		api.POST("/form/:id/responses", formHandler.CreateResponse)
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("server running")
	router.Run(":8080")
}
