package http

import (
	"context"
	"encoding/json"
	"net/http"
	"training1_gokit/middleware"
	"training1_gokit/utils/decodeencode"

	delivery "training1_gokit/modules/account/delivery"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServe create http server with go standard lib
func NewHTTPServe(
	ctx context.Context,
	svcEndpoints delivery.Endpoints,
	logger log.Logger,
) http.Handler {
	// Initialize mux router error logger and error
	var (
		r            = mux.NewRouter()
		options      []httptransport.ServerOption
		errorLogger  = httptransport.ServerErrorLogger(logger)
		errorEncoder = httptransport.ServerErrorEncoder(
			decodeencode.EncodeErrorResponse,
		)
	)
	options = append(options, errorLogger, errorEncoder)

	// Attaching middlewares
	r.Use(middleware.ContentTypeMiddleware)
	r.Use(middleware.AllowOrigin)
	r.Use(mux.CORSMethodMiddleware(r))

	// Creating routes
	r.Methods("POST").Path("/user").Handler(httptransport.NewServer(
		svcEndpoints.CreateUser,
		decodeCreateUserRequest,
		decodeencode.EncodeResponse,
		options...,
	))
	r.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		svcEndpoints.GetUser,
		decodeGetUserRequest,
		decodeencode.EncodeResponse,
		options...,
	))

	return r
}

func decodeCreateUserRequest(
	_ context.Context,
	r *http.Request,
) (interface{}, error) {
	var req delivery.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetUserRequest(
	_ context.Context,
	r *http.Request,
) (interface{}, error) {
	var req delivery.GetUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
