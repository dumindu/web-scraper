package user

import (
	"net/http"

	v "github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	e "web-scraper.dev/internal/api/errors"
	"web-scraper.dev/internal/utils/ctxutil"
	l "web-scraper.dev/internal/utils/logger"
)

type API struct {
	db        *gorm.DB
	logger    *l.Logger
	validator *v.Validate
}

func New(db *gorm.DB, logger *l.Logger, validator *v.Validate) *API {
	return &API{
		db:        db,
		logger:    logger,
		validator: validator,
	}
}

func (a *API) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	form := &FormSignUp{}
	if e.JSONBindAndValidateErrorHandled(w, r, a.logger, a.validator, form, reqID) {
		return
	}
}

func (a *API) SignIn(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	form := &FormSignIn{}
	if e.JSONBindAndValidateErrorHandled(w, r, a.logger, a.validator, form, reqID) {
		return
	}
}
