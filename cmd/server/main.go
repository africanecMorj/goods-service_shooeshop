package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/handler"
	"github.com/africanecMorj/goods-service_shooeshop/internal/middleware"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

func main() {
	log.Println("START APP")

	dbURL := os.Getenv("DATABASE_URL")
	log.Println("DATABASE_URL =", dbURL)

	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("DB init error:", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("DB connected")

	userRepo := &repository.UserRepo{DB: db}
	tokenRepo := &repository.TokenRepo{DB: db}
	productRepo := &repository.ProductRepo{DB: db}
	imageRepo := &repository.ImageRepo{DB: db}

	authService := &service.AuthService{
		Users:  userRepo,
		Tokens: tokenRepo,
		Secret: []byte("secret"),
	}

	productService := &service.ProductService{
		Repo: productRepo,
	}

	streamerService := &service.StreamerService{
		ImageRepo:   imageRepo,
		ProductRepo: productRepo,
	}

	authHandler := &handler.AuthHandler{S: authService}
	productHandler := &handler.ProductHandler{Service: productService}
	streamerHandler := &handler.StreamerHandler{Service: streamerService}

	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth([]byte("secret")))

		r.Route("/products", func(r chi.Router) {

			r.Get("/{id}", productHandler.GetProduct)
			r.Get("/", productHandler.GetProducts)
			r.Get("/{id}/image", streamerHandler.GetImage)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductCreate))
				r.Post("/", productHandler.CreateProduct)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductUpdate))
				r.Patch("/{id}", productHandler.PatchProduct)
				r.Patch("/{id}/image", streamerHandler.PatchImage)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductDelete))
				r.Delete("/{id}", productHandler.DeleteProduct)
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("started :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}