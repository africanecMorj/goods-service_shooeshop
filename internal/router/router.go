package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/handler"
	"github.com/africanecMorj/goods-service_shooeshop/internal/middleware"
	"github.com/redis/go-redis/v9"
)

// Handlers container for DI
type Handlers struct {
	Auth     *handler.AuthHandler
	Product  *handler.ProductHandler
	Streamer *handler.StreamerHandler
	Updater *handler.UpdateHandler
	RDB 	*redis.Client
}

func New(h Handlers) http.Handler {
	r := chi.NewRouter()

	// --- optional middleware ---
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {

		// -------- AUTH --------
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
		})

		r.Route("/password", func(r chi.Router) {
			r.Post("/{code}", h.Updater.ResetPassword)
		
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(h.RDB))
				r.Get("/", h.Updater.RequestPasswordReset)
			})
		})
		// -------- PROTECTED --------
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth([]byte("secret")))

			r.Route("/admin", func(r chi.Router) {
				// ----- READ -----
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequirePermission(domain.UserRead))
					r.Get("/", h.Auth.GetUsers)
				})

				// ----- UPDATE -----
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequirePermission(domain.UserUpdate))
					r.PATCH("/{id}", h.Updater.PatchUser)
				})

			})

			r.Route("/products", func(r chi.Router) {
				// ----- READ -----
				r.Get("/", h.Product.GetProducts)
				r.Get("/{id}", h.Product.GetProduct)
				r.Get("/{id}/image", h.Streamer.GetImage)

				// ----- CREATE -----
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequirePermission(domain.ProductCreate))
					r.Post("/", h.Product.CreateProduct)
				})

				// ----- UPDATE -----
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequirePermission(domain.ProductUpdate))
					r.Patch("/{id}", h.Product.PatchProduct)
					r.Patch("/{id}/image", h.Streamer.PatchImage)
				})

				// ----- DELETE -----
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequirePermission(domain.ProductDelete))
					r.Delete("/{id}", h.Product.DeleteProduct)
				})

			})
		})
	})

	return r
}