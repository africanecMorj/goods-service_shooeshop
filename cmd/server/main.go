package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/handler"
	"github.com/africanecMorj/goods-service_shooeshop/internal/middleware"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service"
)

func main() {
	db, _ := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/shoe_shop")

	userRepo := &repository.UserRepo{DB: db}
	tokenRepo := &repository.TokenRepo{DB: db}
	productRepo := &repository.ProductRepo{DB: db}

	authService := &service.AuthService{
		Users:  userRepo,
		Tokens: tokenRepo,
		Secret: []byte("secret"),
	}

	productService := &service.ProductService{
		Repo: productRepo,
	}

	authHandler := &handler.AuthHandler{S: authService}
	productHandler := &handler.ProductHandler{Service: productService}

	r := chi.NewRouter()

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth([]byte("secret")))

		r.Route("/products", func(r chi.Router) {

			r.Get("/{id}", productHandler.GetProduct)
			r.Get("/", productHandler.GetProducts)
			r.Get("/{id}/image", productHandler.GetImage)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductCreate))
				r.Post("/", productHandler.CreateProduct)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductUpdate))
				r.Patch("/{id}", productHandler.PatchProduct)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductDelete))
				r.Patch("/{id}", productHandler.DeleteProduct)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission(domain.ProductUpdate))
				r.Patch("/{id}/image", productHandler.PatchImage)
			})

		})
	})

	log.Println("started :8080")
	http.ListenAndServe(":8080", r)
}