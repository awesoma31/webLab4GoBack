package common

import (
	"encoding/json"
	"errors"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return fallback
	}
	return v
}

func WriteJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleAndWriteGrpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if ok {
		log.Printf("gRPC error: %s", st.Message())

		var statusCode int
		switch st.Code() {
		case codes.InvalidArgument:
			statusCode = http.StatusBadRequest // 400
		case codes.NotFound:
			statusCode = http.StatusNotFound // 404
		case codes.AlreadyExists:
			statusCode = http.StatusConflict // 409
		default:
			statusCode = http.StatusInternalServerError // 500
		}

		WriteError(w, statusCode, st.Message())
	} else {
		log.Printf("Unexpected error: %v", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be 'Bearer <token>'")
	}

	return parts[1], nil
}

func ReadJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJson(w, status, map[string]string{"error": message})
}
