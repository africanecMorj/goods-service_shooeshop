package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/africanecMorj/goods-service_shooeshop/internal/handler"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/router"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

func main() {
	// --- DB 
	dbURL := os.Getenv("DATABASE_URL")
	dbURL = "postgres://postgres:password@localhost:5432/shoe_shop"

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,  
		Protocol: 2,
	})

	// if dbURL == "" {
	// 	log.Fatal("Empty db_url")
	// }

	db, err := pgxpool.New(context.Background(), dbURL)
	log.Println(err)
	// if err != nil {
	// 	log.Fatal("DB init error:", err)
	// }

	// if err := db.Ping(context.Background()); err != nil {
	// 	log.Fatal("DB ping error:", err)
	// }

	// log.Println("DB connected")


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
		Secret: []byte("secret"),
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