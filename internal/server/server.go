package server

import (
	"net/http"

	"github.com/chukwiz/go-shop/internal/config"
	"github.com/chukwiz/go-shop/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	userService    *services.UserService
	productService *services.ProductService
	uploadService  *services.UploadService
	cartService    *services.CartService
	orderService   *services.OrderService
}

func New(cfg *config.Config,
	db *gorm.DB, logger *zerolog.Logger,
	authService *services.AuthService,
	userService *services.UserService,
	productService *services.ProductService,
	uploadService *services.UploadService,
	cartService *services.CartService,
	orderService *services.OrderService) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		userService:    userService,
		productService: productService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	router.GET("/health", s.healthCheck)
	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authRoutes := auth
			authRoutes.POST("/register", s.register)
			authRoutes.POST("/login", s.login)
			authRoutes.POST("/refresh", s.refreshToken)
			authRoutes.POST("/logout", s.logout)
		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())

		{
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			categories := protected.Group("/categories")
			{
				categoryRoutes := categories
				categoryRoutes.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.CreateProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
			}

			cart := protected.Group("/cart")
			{
				cartRoutes := cart
				cartRoutes.GET("/", s.getCart)
				cartRoutes.POST("/items", s.addToCart)
				cartRoutes.PUT("/items/:id", s.updateCartItem)
				cartRoutes.DELETE("/items/:id", s.removeFromCart)
			}

			orders := protected.Group("/orders")
			{
				ordersRoutes := orders
				ordersRoutes.POST("/", s.createOrder)
				ordersRoutes.GET("/", s.getOrders)
				ordersRoutes.GET("/:id", s.getOrder)
			}

		}
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)
	}
	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}
