package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func NewRouter(handler *AuthHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/api/register", handler.Register)
	r.Post("/api/login", handler.Login)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(jwtSecret))
		r.Get("/api/dashboard/user", handler.GetUserDashboard)
		r.Group(func(r chi.Router) {
			r.Use(AdminOnlyMiddleware)
			r.Get("/api/dashboard/admin", handler.GetAdminDashboard)
		})
	})
	return r
}