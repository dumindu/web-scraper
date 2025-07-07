package errors

import (
	"net/http"
)

var (
	RespDBDataInsertFailure = []byte(`{"error": "db data insert failure"}`)
	RespDBDataAccessFailure = []byte(`{"error": "db data access failure"}`)
	RespDBDataUpdateFailure = []byte(`{"error": "db data update failure"}`)
	RespDBDataDeleteFailure = []byte(`{"error": "db data delete failure"}`)

	RespJSONEncodeFailure = []byte(`{"error": "json encode failure"}`)
	RespJSONDecodeFailure = []byte(`{"error": "json decode failure"}`)

	RespHashGenerationFailure = []byte(`{"error": "hash generation failure"}`)
	RespEmailSendingFailure   = []byte(`{"error": "email sending failure"}`)

	RespInvalidActivationRequest = []byte(`{"error": "invalid activation request"}`)
	RespTokenExpired             = []byte(`{"error": "token expired"}`)
	RespTokenInvalid             = []byte(`{"error": "invalid token"}`)
	RespUnauthorized             = []byte(`{"error": "unauthorized"}`)
)

type Error struct {
	Error string `json:"error"`
}

type ValidationErrors struct {
	Errors []string `json:"errors"`
}

func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(error)
}

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error)
}

func Unauthorized(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(error)
}

func Conflict(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusConflict)
	w.Write(error)
}
