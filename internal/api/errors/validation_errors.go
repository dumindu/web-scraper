package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	l "web-scraper.dev/internal/utils/logger"
)

func JSONBindAndValidateErrorHandled(w http.ResponseWriter, r *http.Request, logger *l.Logger, validator *validator.Validate, form any, reqID string) bool {
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		BadRequest(w, RespJSONDecodeFailure)

		return true
	}

	if err := validator.Struct(form); err != nil {
		respBody, err := json.Marshal(ToValidationErrors(err))
		if err != nil {
			logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			ServerError(w, RespJSONEncodeFailure)

			return true
		}

		BadRequest(w, respBody)
		return true
	}

	return false
}

func ToValidationErrors(err error) *ValidationErrors {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		resp := ValidationErrors{
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "eqfield":
				resp.Errors[i] = fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
			case "min":
				resp.Errors[i] = fmt.Sprintf("%s must be at least in %v characters", err.Field(), err.Param())
			case "email":
				resp.Errors[i] = fmt.Sprintf("%s must be a valid email address", err.Field())
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}
