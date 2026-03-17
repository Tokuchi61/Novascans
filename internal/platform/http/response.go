package http

import (
	"encoding/json"
	stdhttp "net/http"
)

type dataEnvelope struct {
	Data any `json:"data"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteData(w stdhttp.ResponseWriter, status int, data any) {
	WriteJSON(w, status, dataEnvelope{Data: data})
}

func WriteNoContent(w stdhttp.ResponseWriter) {
	w.WriteHeader(stdhttp.StatusNoContent)
}

func WriteError(w stdhttp.ResponseWriter, err error) {
	appErr := ResolveError(err)

	WriteJSON(w, appErr.StatusCode, errorEnvelope{
		Error: errorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

func WriteJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusInternalServerError), stdhttp.StatusInternalServerError)
	}
}
