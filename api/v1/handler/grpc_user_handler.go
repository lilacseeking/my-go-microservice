// api/v1/handler/grpc_user_handler.go
package handler

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"my-go-microservice/internal/model"
	"my-go-microservice/internal/service"
	pb "my-go-microservice/proto"
)

// GRPCUserHandler 实现 gRPC UserServiceServer 接口
type GRPCUserHandler struct {
	pb.UnimplementedUserServiceServer
	UserService service.UserService
}

// NewGRPCUserServer 创建 gRPC 用户服务处理器
func NewGRPCUserServer(userService service.UserService) *GRPCUserHandler {
	return &GRPCUserHandler{
		UserService: userService,
	}
}

// CreateUser 处理创建用户的 gRPC 请求
func (h *GRPCUserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	logrus.WithField("email", req.User.Email).Info("gRPC: 收到创建用户请求")

	// 将 proto User 转换为 model User
	user := h.protoToModel(req.User)

	// 调用业务层创建用户
	if err := h.UserService.CreateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("gRPC: 创建用户失败")
		return nil, err
	}

	logrus.WithField("user_id", user.ID).Info("gRPC: 用户创建成功")

	return &pb.CreateUserResponse{
		UserId: user.ID,
	}, nil
}

// GetUser 处理获取用户的 gRPC 请求
func (h *GRPCUserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	logrus.WithField("user_id", req.Id).Info("gRPC: 收到获取用户请求")

	// 调用业务层获取用户
	user, err := h.UserService.GetUser(ctx, req.Id)
	if err != nil {
		logrus.WithError(err).WithField("user_id", req.Id).Warn("gRPC: 用户不存在")
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	logrus.WithField("user_id", user.ID).Info("gRPC: 用户查询成功")

	// 将 model User 转换为 proto User
	protoUser := h.modelToProto(user)

	return &pb.GetUserResponse{
		User: protoUser,
	}, nil
}

// protoToModel 将 proto User 转换为 model User
func (h *GRPCUserHandler) protoToModel(protoUser *pb.User) *model.User {
	return &model.User{
		ID:    protoUser.Id,
		Name:  protoUser.Name,
		Email: protoUser.Email,
	}
}

// modelToProto 将 model User 转换为 proto User
func (h *GRPCUserHandler) modelToProto(modelUser *model.User) *pb.User {
	return &pb.User{
		Id:        modelUser.ID,
		Name:      modelUser.Name,
		Email:     modelUser.Email,
		CreatedAt: timestamppb.New(modelUser.CreatedAt),
	}
}
