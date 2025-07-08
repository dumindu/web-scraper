package keyword

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	v "github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	e "web-scraper.dev/internal/api/errors"
	"web-scraper.dev/internal/repository"
	"web-scraper.dev/internal/utils/ctxutil"
	l "web-scraper.dev/internal/utils/logger"

	_ "web-scraper.dev/internal/model"
)

type API struct {
	db        *repository.Db
	logger    *l.Logger
	validator *v.Validate
}

func New(db *gorm.DB, logger *l.Logger, validator *v.Validate) *API {
	return &API{
		db:        repository.New(db),
		logger:    logger,
		validator: validator,
	}
}

// GetKeywords godoc
// @summary Get the list of keywords
// @description Get the list of keywords uploaded by current user
// @tags keywords
//
// @router /keywords [GET]
// @accept json
// @produce json
// @security BearerToken
// @success 200 {array} model.KeywordDTO
// @failure 401 {object} e.Error
// @failure 500 {object} e.Error
func (a *API) GetKeywords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, ctxUser := ctxutil.RequestID(ctx), ctxutil.UserFromCtx(ctx)

	keywords, err := a.db.ListKeywordsByUserId(*ctxUser.ID)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if keywords == nil {
		fmt.Fprint(w, "[]")
		return
	}

	dto := keywords.ToDTOs()
	if err := json.NewEncoder(w).Encode(&dto); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// GetKeyword godoc
// @summary Get the result of a keyword
// @description Get the result of a keyword uploaded by current user
// @tags keywords
//
// @router /keywords/{id} [GET]
// @accept json
// @produce json
// @security BearerToken
// @param id path string true "Payment ID in uuid format"
//
// @success 200 {object} model.KeywordDTO
// @failure 400 {object} e.Error
// @failure 401 {object} e.Error
// @failure 404
// @failure 500 {object} e.Error
func (a *API) GetKeyword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, ctxUser := ctxutil.RequestID(ctx), ctxutil.UserFromCtx(ctx)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		e.BadRequest(w, e.RespInvalidID)
	}

	keyword, err := a.db.ReadKeywordByIdAndUserId(int64(id), *ctxUser.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(&keyword); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

func (a *API) UploadKeywords(w http.ResponseWriter, r *http.Request) {
}
