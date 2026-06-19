package responseHTTP

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

func JSONError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	resp := ErrorResponse{Message: errMessage, Code: code}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Debug("Помилка у кодуванні JSONError:", "err", err.Error())
	}
}

func JSONResp(w http.ResponseWriter, payload any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		slog.Debug("Помилка у кодуванні JSONResp:", "err", err.Error())
	}
}

func JSONRespMessage(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	resp := ErrorResponse{Message: message, Code: code}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Debug("Помилка у кодуванні JSONRespMessage:", "err", err.Error())
	}
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		// slog.Debug("Помилка декодування JSON тіла запиту:", "error", err.Error())
		return err
	}

	return nil
}
