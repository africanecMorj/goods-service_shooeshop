package main

import (
	"log"
	"net/http"
	"os"

	"github.com/africanecMorj/goods-service_shooeshop/internal/handler"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/router"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
	"github.com/africanecMorj/goods-service_shooeshop/cmd/config"
)

func main() {
	// --- DB 
	rdb := config.RedisInit()
	db := config.PostgresInit()

	// -------- REPOSITORIES --------
	userRepo := &repository.UserRepo{DB: db}
	tokenRepo := &repository.TokenRepo{DB: db}
	productRepo := &repository.ProductRepo{DB: db}
	imageRepo := &repository.ImageRepo{DB: db}
	chacheRepo := &repository.ChacheRepo{RDB: rdb}

	// -------- SERVICES --------
	authService := &service.AuthService{
		Users:  userRepo,
		Tokens: tokenRepo,
		Secret: os.Getenv("JWT-SECRET"),
	}

	productService := &service.ProductService{
		Repo: productRepo,
		DB: db,
	}

	streamerService := &service.StreamerService{
		ImageRepo:   imageRepo,
		ProductRepo: productRepo,
		DB: db,
	}

	updateService := &service.UpdateService{
		UserRepo:   userRepo,
		ChacheRepo: chacheRepo,
	}

	// -------- HANDLERS --------
	authHandler := &handler.AuthHandler{S: authService}
	productHandler := &handler.ProductHandler{Service: productService,}
	streamerHandler := &handler.StreamerHandler{Service: streamerService,}
	updateHandler := &handler.UpdateHandler{UpdateService: updateService, AuthService: authService}

	// -------- ROUTER --------
	r := router.New(router.Handlers{
		Auth:     authHandler,
		Product:  productHandler,
		Streamer: streamerHandler,
		Updater: updateHandler,
		RDB:rdb,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}