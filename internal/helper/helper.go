package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
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

func ReadParam(r *http.Request, key string) string {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	return id
}

func DateTime(datetime time.Time) string {
	loc, err := time.LoadLocation("America/Sao_Paulo")

	if err != nil {
		slog.Error("Unable to load timezone", "error", err.Error())
		return ""
	}
	return datetime.In(loc).Format("2006-01-02T15:04:05")
}

func GenerateUUIDV7() string {
	value, err := uuid.NewV7()
	if err != nil {
		slog.Error("Unable to generate UUIDv7", "error", err.Error())
		return ""
	}

	return value.String()
}
