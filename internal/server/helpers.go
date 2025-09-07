package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func respondStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func respondJSON(ctx context.Context, w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(response); err != nil {
		slog.ErrorContext(ctx, "failed to write response", "err", err)
	}
}

type errorMessage struct {
	Error string `json:"error"`
}

func respondError(ctx context.Context, w http.ResponseWriter, code int, message string) {
	respondJSON(ctx, w, code, errorMessage{message})
}

type dataMessage struct {
	Data any `json:"data"`
}

func respondData(ctx context.Context, w http.ResponseWriter, code int, payload any) {
	respondJSON(ctx, w, code, dataMessage{payload})
}

type validator interface {
	Validate() error
}

func bindJSON(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	if v, ok := payload.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validate: %w", err)
		}
	}

	return nil
}

func queryOr[T any](r *http.Request, key string, bo T) T {
	val := r.URL.Query().Get(key)
	if val == "" {
		return bo
	}

	switch any(bo).(type) {
	case string:
		return any(val).(T)
	case int:
		if intVal, err := strconv.Atoi(val); err == nil {
			return any(intVal).(T)
		}
	case uint64:
		if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
			return any(uintVal).(T)
		}
	}

	return bo
}
