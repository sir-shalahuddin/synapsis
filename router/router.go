package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/config"
	"github.com/sir-shalahuddin/synapsis/controller"
	"github.com/sir-shalahuddin/synapsis/middleware"
	"github.com/sir-shalahuddin/synapsis/pkg/auth"
	"github.com/sir-shalahuddin/synapsis/repository"
	"github.com/sir-shalahuddin/synapsis/service"
)

func SetupRoutes(app *fiber.App, db *sqlx.DB, config config.JWTConfig) {
	userRepo := repository.NewUserRepository(db)
	tokenManager := auth.NewManager(config.Secret)
	userService := service.NewUserService(userRepo, tokenManager)
	userController := controller.NewUserController(userService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewproductService(productRepo)
	productController := controller.NewProductController(productService)

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartController := controller.NewCartController(cartService)

	orderRepo := repository.NewOrderRepository(db)
	txRepo := repository.NewTxRepository(db)
	orderService := service.NewOrderService(orderRepo, txRepo)
	orderController := controller.NewOrderController(orderService)

	// middleware logger
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", userController.Login)
	auth.Post("/register", userController.Register)

	// Product
	product := api.Group("/products")
	product.Get("/", productController.GetProducts)

	// Cart
	cart := api.Group("/carts", middleware.Protected(config))
	cart.Post("/", cartController.AddProduct)
	cart.Get("/", cartController.GetProducts)
	cart.Delete("/", cartController.DeleteProduct)

	// Order
	order := api.Group("/orders", middleware.Protected(config))
	order.Post("/", orderController.CreateOrder)
	order.Post("/:order_id/payments", orderController.PayOrder)

}
