package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/cookie"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
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

	log := logger.NewLogger()
	log.SetComponentName("helper")
	log.SetMethodName("BadRequest")

	resp := response.ClientError{
		Status:    400,
		Code:      "BAD_REQUEST",
		Message:   err.Error(),
		Timestamp: DateTime(time.Now().UTC()),
	}

	log.Error("bad request", "response", resp)

	WriteJSON(w, http.StatusBadRequest, map[string]any{"error": resp})
}

func NotFound(w http.ResponseWriter, err error) {
	log := logger.NewLogger()
	log.SetComponentName("helper")
	log.SetMethodName("NotFound")
	log.Error("resource not found", "error", err.Error())

	empty := map[string]any{}

	WriteJSON(w, http.StatusNotFound, map[string]any{"error": "resource not found", "data": empty})
}

func InternalServerError(w http.ResponseWriter, err error) {
	log := logger.NewLogger()
	log.SetComponentName("helper")
	log.SetMethodName("InternalServerError")

	empty := map[string]any{}

	log.Error("internal server error", "error", err.Error())
	WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "Internal Server Error", "data": empty})
}

func ReadParam(r *http.Request, key string) string {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	return id
}

func DateTime(datetime time.Time) string {
	loc, err := time.LoadLocation("America/Sao_Paulo")

	if datetime.IsZero() {
		return ""
	}

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

func SetCookie(w http.ResponseWriter, ac cookie.Account) {
	json, err := json.Marshal(ac)
	if err != nil {
		log.Fatalf("error marshaling JSON: %v", err)
	}

	content := base64.StdEncoding.EncodeToString(json)

	cookie := http.Cookie{
		Name:     "account",
		Value:    content,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetCookie(w http.ResponseWriter, r *http.Request) *cookie.Account {
	log := logger.NewLogger()
	log.SetComponentName("Helper")
	log.SetMethodName("GetCookie")

	c, err := r.Cookie("account")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Error("cookie not found")
			BadRequest(w, err)
		default:
			log.Error("generic error", err.Error())
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return nil
	}

	var ac *cookie.Account

	content, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		log.Error("error cooking decode base64", "error", err.Error())
	}

	err = json.Unmarshal([]byte(content), &ac)
	if err != nil {
		log.Error("error unmarshaling JSON", "json", c.Value)
	}
	return ac
}
