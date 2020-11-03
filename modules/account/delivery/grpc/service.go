package grpc

import (
	"context"

	delivery "training1_gokit/modules/account/delivery"
	"training1_gokit/modules/account/delivery/protobuf/account_grpc"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	createuser kitgrpc.Handler
	getuser    kitgrpc.Handler
	logger     log.Logger
}

// NewGRPCServer create grpc server
func NewGRPCServer(
	svcEndpoints delivery.Endpoints,
	logger log.Logger,
) account_grpc.AccountServiceServer {
	var options []kitgrpc.ServerOption
	errorLogger := kitgrpc.ServerErrorLogger(logger)
	options = append(options, errorLogger)

	return &grpcServer{
		createuser: kitgrpc.NewServer(
			svcEndpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
			options...,
		),
		getuser: kitgrpc.NewServer(
			svcEndpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
			options...,
		),
		logger: logger,
	}
}

func (s *grpcServer) CreateUser(
	ctx oldcontext.Context, req *account_grpc.CreateUserRequest,
) (*account_grpc.CreateUserResponse, error) {
	_, rep, err := s.createuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*account_grpc.CreateUserResponse), nil
}

func (s *grpcServer) GetUser(
	ctx oldcontext.Context, req *account_grpc.GetUserRequest,
) (*account_grpc.GetUserResponse, error) {
	_, rep, err := s.getuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*account_grpc.GetUserResponse), nil
}

// decodeCreateUserRequest to json
func decodeCreateUserRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req := request.(*account_grpc.CreateUserRequest)
	return delivery.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

// decodeGetUserRequest to json
func decodeGetUserRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req := request.(*account_grpc.GetUserRequest)
	return delivery.GetUserRequest{
		Id: req.Id,
	}, nil
}

// encodeCreateUserResponse to json
func encodeCreateUserResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	res := response.(delivery.CreateUserResponse)
	return &account_grpc.CreateUserResponse{Ok: res.Ok}, nil
}

// encodeGetUserResponse to json
func encodeGetUserResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	res := response.(delivery.GetUserResponse)
	return &account_grpc.GetUserResponse{Email: res.Email}, nil
}
