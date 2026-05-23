package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, dst any) error {

	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		slog.Error("Error when parsing the JSON Body")
		return err
	}

	return nil
}
