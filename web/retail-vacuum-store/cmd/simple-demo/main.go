package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/handlers"
	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	zerolog_log "github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	setupLogging()

	// Setup mock services for demo (no database required)
	logger := zerolog_log.With().Str("component", "api").Logger()
	productService := services.NewMockProductService(logger)
	cartService := services.NewMockCartService(logger)

	// Setup handlers
	productHandler := handlers.NewProductHandler(productService, logger)
	cartHandler := handlers.NewCartHandler(cartService, logger)

	// Setup Echo server
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Routes
	setupRoutes(e, productHandler, cartHandler)

	// Start server
	port := "3001"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	
	serverAddr := fmt.Sprintf(":%s", port)
	zerolog_log.Info().Str("address", serverAddr).Msg("Starting demo server (mock data)")

	// Start server in a goroutine
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			zerolog_log.Error().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	zerolog_log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		zerolog_log.Error().Err(err).Msg("Server forced to shutdown")
	}

	zerolog_log.Info().Msg("Server shutdown complete")
}

func setupLogging() {
	// Setup structured logging
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use console writer for development
	zerolog_log.Logger = zerolog_log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupRoutes(e *echo.Echo, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler) {
	// API version 1
	api := e.Group("/api/v1")

	// Product routes
	products := api.Group("/products")
	products.GET("", productHandler.SearchProducts)
	products.POST("", productHandler.CreateProduct)
	products.GET("/featured", productHandler.GetFeaturedProducts)
	products.GET("/categories", productHandler.GetCategories)
	products.GET("/brands", productHandler.GetBrands)
	products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
	products.GET("/brand/:brandId", productHandler.GetProductsByBrand)
	products.GET("/:id", productHandler.GetProduct)
	products.PUT("/:id", productHandler.UpdateProduct)
	products.DELETE("/:id", productHandler.DeleteProduct)
	products.GET("/sku/:sku", productHandler.GetProductBySKU)

	// Cart routes
	cart := api.Group("/cart")
	cart.GET("", cartHandler.GetCart)
	cart.POST("/items", cartHandler.AddToCart)
	cart.PUT("/items/:productId", cartHandler.UpdateCartItem)
	cart.DELETE("/items/:productId", cartHandler.RemoveFromCart)
	cart.DELETE("", cartHandler.ClearCart)
	cart.GET("/summary", cartHandler.GetCartSummary)
	cart.POST("/reserve", cartHandler.ReserveCartItems)
	cart.POST("/release", cartHandler.ReleaseCartReservations)
	cart.POST("/merge", cartHandler.MergeGuestCart)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
			"mode":   "demo",
		})
	})

	// API documentation endpoint
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"name":    "Vacuum Cleaner Retail Store API (Demo Mode)",
			"version": "1.0.0",
			"mode":    "demo",
			"endpoints": map[string]interface{}{
				"products": map[string]string{
					"GET /api/v1/products":                    "Search products with filters",
					"POST /api/v1/products":                   "Create a new product",
					"GET /api/v1/products/featured":           "Get featured products",
					"GET /api/v1/products/categories":         "Get all categories",
					"GET /api/v1/products/brands":             "Get all brands",
					"GET /api/v1/products/category/:id":       "Get products by category",
					"GET /api/v1/products/brand/:id":          "Get products by brand",
					"GET /api/v1/products/:id":                "Get product by ID",
					"PUT /api/v1/products/:id":                "Update product",
					"DELETE /api/v1/products/:id":             "Delete product",
					"GET /api/v1/products/sku/:sku":           "Get product by SKU",
				},
				"cart": map[string]string{
					"GET /api/v1/cart":                        "Get current cart",
					"POST /api/v1/cart/items":                 "Add item to cart",
					"PUT /api/v1/cart/items/:productId":       "Update cart item quantity",
					"DELETE /api/v1/cart/items/:productId":    "Remove item from cart",
					"DELETE /api/v1/cart":                     "Clear cart",
					"GET /api/v1/cart/summary":                "Get cart summary",
					"POST /api/v1/cart/reserve":               "Reserve cart items",
					"POST /api/v1/cart/release":               "Release cart reservations",
					"POST /api/v1/cart/merge":                 "Merge guest cart with customer cart",
				},
				"health": map[string]string{
					"GET /health": "Health check",
				},
			},
		})
	})
	
	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Welcome to the Vacuum Cleaner Retail Store API",
			"docs":    "/api",
			"health":  "/health",
		})
	})
}
