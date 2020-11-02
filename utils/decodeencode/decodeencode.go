package decodeencode

import (
	"context"
	"encoding/json"
	"net/http"
)

// Custom error type for business logic error
type errorer interface {
	error() error
}

// Identify error and returns http error code
func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

// EncodeErrorResponse error response decoder for all services
func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// EncodeResponse encode response for all services
func EncodeResponse(
	ctx context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		EncodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
