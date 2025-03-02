package handler

import (
	"context"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	v1 "github.com/iliadmitriev/go-user-test/internal/server/grpc/user/v1"
	"github.com/iliadmitriev/go-user-test/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc -I=../../grpc/ --go_out=../../internal/server/grpc --go_opt=paths=source_relative --go-grpc_out=../../internal/server/grpc --go-grpc_opt=paths=source_relative ../../grpc/user/v1/user.proto

type GRPCHandler interface {
	RegisterGRPC(srv *grpc.Server)
}

type grpcUserHandler struct {
	v1.UserServiceServer

	userService service.UserServiceInterface
	logger      *zap.SugaredLogger
}

func NewGRPCUserHandler(userService service.UserServiceInterface, logger *zap.Logger) GRPCHandler {
	return &grpcUserHandler{
		userService: userService,
		logger:      logger.Named("GRPCUserHandler").Sugar(),
	}
}

func (g *grpcUserHandler) RegisterGRPC(srv *grpc.Server) {
	v1.RegisterUserServiceServer(srv, g)
}

func (g *grpcUserHandler) Create(ctx context.Context, r *v1.CreateRequest) (*v1.CreateResponse, error) {
	var userIn domain.UserIn
	userIn.Login = r.GetLogin()
	userIn.Password = r.GetPassword()
	userIn.Name = r.GetName()

	g.logger.Infow("Got grpc request", "login", userIn.Login, "name", userIn.Name)
	_, err := g.userService.CreateUser(ctx, &userIn)
	if err != nil {
		g.logger.Warnw("Error creating user", "err", err)
		return nil, err
	}

	var userOutGrpc v1.CreateResponse
	userOutGrpc.Code = 200
	userOutGrpc.Message = "OK"
	return &userOutGrpc, nil
}

func (g *grpcUserHandler) GetByLogin(ctx context.Context, r *v1.GetByLoginRequest) (*v1.GetUserResponse, error) {
	user, err := g.userService.GetUser(ctx, r.GetLogin())
	if err != nil {
		g.logger.Warnw("Error getting user", "err", err)
		return nil, err
	}

	id, err := user.ID.MarshalBinary()
	if err != nil {
		g.logger.Warnw("Error marshaling user id", "err", err)
		return nil, err
	}

	return &v1.GetUserResponse{
		Id:        id,
		Login:     user.Login,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
