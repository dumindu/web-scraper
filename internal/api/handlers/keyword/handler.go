package keyword

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	v "github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	e "web-scraper.dev/internal/api/errors"
	"web-scraper.dev/internal/model"
	"web-scraper.dev/internal/repository"
	"web-scraper.dev/internal/tasks"
	"web-scraper.dev/internal/utils/ctxutil"
	l "web-scraper.dev/internal/utils/logger"
)

const SearchEngineBing = "bing"

type API struct {
	db        *repository.Db
	logger    *l.Logger
	validator *v.Validate
	asyq      *asynq.Client
}

func New(db *gorm.DB, logger *l.Logger, validator *v.Validate, asyq *asynq.Client) *API {
	return &API{
		db:        repository.New(db),
		logger:    logger,
		validator: validator,
		asyq:      asyq,
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
// @param id path string true "Keyword ID in uuid format"
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
		return
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

// UploadKeywords godoc
// @summary Upload the keywords CSV file
// @description Upload the keywords CSV file to scrape on web
// @tags keywords
//
// @router /keywords [POST]
// @Accept multipart/form-data
// @Accept text/csv
// @produce json
// @security BearerToken
// @Param file formData file true "CSV file with keywords"
//
// @success 202
// @failure 400 {object} e.Error
// @failure 401 {object} e.Error
// @failure 500 {object} e.Error
func (a *API) UploadKeywords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, ctxUser := ctxutil.RequestID(ctx), ctxutil.UserFromCtx(ctx)

	file, _, err := r.FormFile("file")
	if err != nil {
		e.BadRequest(w, e.RespInvalidFile)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil || len(records) == 0 {
		e.BadRequest(w, e.RespInvalidFile)
		return
	}

	var keywords []string

	for _, record := range records {
		if len(record) > 0 && strings.TrimSpace(record[0]) != "" {
			keywords = append(keywords, strings.TrimSpace(record[0]))
		}

		if len(keywords) > 100 {
			e.BadRequest(w, e.RespInvalidFileExceedMaxRows)
			return
		}
	}

	if len(keywords) == 0 {
		e.BadRequest(w, e.RespInvalidFile)
		return
	}

	userID := *ctxUser.ID
	keywordModels := make([]*model.Keyword, len(keywords))
	for i, v := range keywords {
		keywordModels[i] = &model.Keyword{
			UserID:       userID,
			Keyword:      v,
			Status:       "pending",
			SearchEngine: SearchEngineBing,
		}
	}

	if err := a.db.Create(&keywordModels).Error; err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	for _, v := range keywordModels {
		task := tasks.NewScrapeKeywordTask(v.ID)
		// Enqueue with a delay to avoid rate limiting
		if _, err := a.asyq.Enqueue(task, asynq.ProcessIn(tasks.ScrapeKeywordDelayInSeconds*time.Second)); err != nil {
			a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Str("task", "scrape-keyword").Msg("")
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
