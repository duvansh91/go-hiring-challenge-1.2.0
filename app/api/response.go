package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("error parsing response", "error", err)
		http.Error(w, "error parsing response", http.StatusInternalServerError)

		return
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		slog.Error("failed parsing error response", "error", err)
		http.Error(w, "failed parsing error response", http.StatusInternalServerError)

		return
	}
}
