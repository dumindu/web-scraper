package user

import (
	"errors"
	"net/http"
	"strings"

	v "github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	e "web-scraper.dev/internal/api/errors"
	"web-scraper.dev/internal/mailer"
	"web-scraper.dev/internal/model"
	"web-scraper.dev/internal/repository"
	"web-scraper.dev/internal/utils/ctxutil"
	l "web-scraper.dev/internal/utils/logger"
)

type API struct {
	db        *repository.Db
	mailer    *mailer.Mailer
	logger    *l.Logger
	validator *v.Validate
}

func New(db *gorm.DB, mailer *mailer.Mailer, logger *l.Logger, validator *v.Validate) *API {
	return &API{
		db:        repository.New(db),
		mailer:    mailer,
		logger:    logger,
		validator: validator,
	}
}

// SignUp godoc
// @summary User signup
// @description User signup by using the given email and password.
// @description This will send the email verification email with a 6-character long code to the given email address.
//
// @tags users
//
// @router /users/sign-up [POST]
// @accept json
// @produce  json
// @param body body FormSignUp true "SignUp Form"
//
// @success 201 "Created"
// @success 409 "Conflict"
// @failure 400 {object} e.Error
// @failure 400 {object} e.ValidationErrors
// @failure 500 {object} e.Error
func (a *API) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	form := &FormSignUp{}
	if e.JSONBindAndValidateErrorHandled(w, r, a.logger, a.validator, form, reqID) {
		return
	}

	form.Email = strings.ToLower(form.Email)

	user, err := a.db.ReadUserByEmail(form.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if user != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespHashGenerationFailure)
		return
	}

	userModel := model.NewUser(form.Email, string(hashedPassword))
	if err := a.db.CreateUser(userModel); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	if err := a.mailer.ActivationMail(form.Email, userModel.ActivationToken.Token); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespEmailSendingFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) SignIn(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	form := &FormSignIn{}
	if e.JSONBindAndValidateErrorHandled(w, r, a.logger, a.validator, form, reqID) {
		return
	}
}
