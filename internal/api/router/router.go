package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"web-scraper.dev/internal/api/handlers/health"
	"web-scraper.dev/internal/api/handlers/keyword"
	"web-scraper.dev/internal/api/handlers/user"
	"web-scraper.dev/internal/api/router/middleware"
	"web-scraper.dev/internal/api/router/middleware/requestlog"
	"web-scraper.dev/internal/mailer"
	"web-scraper.dev/internal/utils/logger"
)

func New(hd time.Duration, hdw time.Duration, db *gorm.DB, ml *mailer.Mailer, l *logger.Logger, v *validator.Validate, asyq *asynq.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "pragma"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-Total-Count"},
		MaxAge:           300,
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(cors.Handler)
		r.Use(middleware.ContentTypeJSON)
		r.Use(middleware.RequestID)

		userAPI := user.New(db, ml, l, v)
		r.Method(http.MethodPost, "/users/sign-up", requestlog.NewHandler(userAPI.SignUp, hd, l))
		r.Method(http.MethodPost, "/users/sign-in", requestlog.NewHandler(userAPI.SignIn, hd, l))

		r.Method(http.MethodPost, "/users/activate", requestlog.NewHandler(userAPI.Activate, hd, l))

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.JwtAuthentication)

			keywordAPI := keyword.New(db, l, v, asyq)
			r.Method(http.MethodGet, "/keywords", requestlog.NewHandler(keywordAPI.GetKeywords, hd, l))
			r.Method(http.MethodGet, "/keywords/{id}", requestlog.NewHandler(keywordAPI.GetKeyword, hd, l))

			r.Method(http.MethodPost, "/keywords", requestlog.NewHandler(keywordAPI.UploadKeywords, hd, l))
		})
	})

	return r
}
