package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Respond sends a response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	setStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}
