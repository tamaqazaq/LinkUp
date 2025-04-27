package main

import (
	"WebMessanger/internal/adapter/handler"
	"WebMessanger/internal/adapter/postgres"
	"WebMessanger/internal/app/service"
	"WebMessanger/pkg/middleware"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := postgres.NewUserRepo(db)
	userService := service.NewUserService(repo)
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	threadRepo := postgres.NewThreadRepo(db)
	threadService := service.NewThreadService(threadRepo)
	threadHandler := handler.NewThreadHandler(threadService)
	commentRepo := postgres.NewCommentRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService, userService)
	likeRepo := postgres.NewLikeRepo(db)
	likeService := service.NewLikeService(likeRepo)
	likeHandler := handler.NewLikeHandler(likeService, userService, threadService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://linkup-9w5.pages.dev",
			"https://*.linkup-9w5.pages.dev",
		},
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://linkup-9w5.pages.dev" ||
				len(origin) > 0 &&
					origin[len(origin)-19:] == ".linkup-9w5.pages.dev"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/search", userHandler.SearchUsers)
	r.GET("/verify", authHandler.VerifyEmail)

	protected := r.Group("/user", middleware.JWTMiddleware())
	{
		protected.GET("/:id", userHandler.GetUserByID)
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/me", userHandler.UpdateProfile)
		protected.POST("/threads", threadHandler.Create)
		protected.PUT("/threads/:id", threadHandler.Update)
		protected.DELETE("/threads/:id", threadHandler.Delete)
		protected.GET("/threads", threadHandler.GetAll)
		protected.GET("/threads/:id", threadHandler.GetById)
		protected.GET("/users/:user_id/threads", threadHandler.GetByUser)
		protected.GET("/users/:user_id/posts", threadHandler.GetByUser)
	}
	commentRoutes := protected.Group("/comments")
	{
		commentRoutes.POST("/:thread_id", commentHandler.Create)
		commentRoutes.GET("/:thread_id", commentHandler.GetByThread)
		commentRoutes.DELETE("/:comment_id", commentHandler.Delete)
	}
	likeRoutes := protected.Group("/likes")
	{
		likeRoutes.POST("/:thread_id", likeHandler.Create)
		likeRoutes.DELETE("/:thread_id", likeHandler.RemoveLike)
		likeRoutes.GET("/:thread_id", likeHandler.GetLikesByThread)
		likeRoutes.GET("/user/:user_id", likeHandler.GetLikesByUser)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}
