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

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		slog.Error("Error when parsing the JSON Body")
		return err
	}

	return nil
}

func BadRequest(w http.ResponseWriter, err error) {
	slog.Error("Bad request: %+v", err)
	WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
}

func InternalServerError(w http.ResponseWriter, err error) {
	slog.Error("Internal server error: %+v", err)
	WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
}
