package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorCode string

type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func Error(w http.ResponseWriter, statusCode int, code ErrorCode, msg string, args ...any) {
	details := argsToMap(args)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"error": ErrorResponse{
			Code:    string(code),
			Message: msg,
			Details: details,
		},
	})
}

func argsToMap(args []any) map[string]any {
	if len(args) <= 0 {
		return nil
	}

	if len(args)%2 != 0 {
		return map[string]any{
			"args": args,
		}
	}

	m := make(map[string]any)
	for i := range len(args) / 2 {
		key := args[i*2+0].(string)
		value := fmt.Sprintf("%v", args[i*2+1])
		m[key] = value
	}

	return m
}

func ErrorBadRequest(w http.ResponseWriter, msg string, args ...any) {
	Error(w, http.StatusBadRequest, "invalid_input", msg, args...)
}

func ErrorInternal(w http.ResponseWriter, msg string, args ...any) {
	Error(w, http.StatusInternalServerError, "internal_error", msg, args...)
}

func OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
