package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func JSONResponse(w http.ResponseWriter, status int, data any) {
	JSON(w, status, data)
}

func JSONErrorResponse(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]any{"error": msg})
}

func JSONValidationErrorResponse(w http.ResponseWriter, fields map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, map[string]any{
		"error":  "validation_failed",
		"fields": fields,
	})
}
