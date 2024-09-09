package handlers

import (
	"encoding/json"
	"net/http"

	"tincho.dev/rest-ws/dto"
	"tincho.dev/rest-ws/server"
)

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := dto.HomeResponse{
			Message: "Welcome to the API",
			Status:  true,
		}

		json.NewEncoder(w).Encode(response)
	}
}
