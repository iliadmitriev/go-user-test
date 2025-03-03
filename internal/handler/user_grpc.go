package handler

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	user_proto "github.com/iliadmitriev/go-user-test/internal/server/grpc/user/v1"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

//go:generate protoc -I=../../grpc/ --go_out=../../internal/server/grpc --go_opt=paths=source_relative --go-grpc_out=../../internal/server/grpc --go-grpc_opt=paths=source_relative ../../grpc/user/v1/user.proto

type GRPCHandler interface {
	RegisterGRPC(srv *grpc.Server)
}

type grpcUserHandler struct {
	user_proto.UserServiceServer

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
	user_proto.RegisterUserServiceServer(srv, g)
}

func (g *grpcUserHandler) Create(ctx context.Context, r *user_proto.CreateRequest) (*user_proto.CreateResponse, error) {
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

	return &user_proto.CreateResponse{
		Code:    200,
		Message: "User created successfully",
	}, nil
}

func (g *grpcUserHandler) GetByLogin(ctx context.Context, r *user_proto.GetByLoginRequest) (*user_proto.GetUserResponse, error) {
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

	return &user_proto.GetUserResponse{
		Id:        id,
		Login:     user.Login,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
