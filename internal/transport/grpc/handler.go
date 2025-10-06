package grpc

import (
	"context"
	"errors"
	"fmt"

	taskpb "github.com/guarref/project-protos/proto/task"
	userpb "github.com/guarref/project-protos/proto/user"
	"github.com/guarref/tasks-service/internal/task"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
  svc        *task.TaskService
  userClient userpb.UserServiceClient
  taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.TaskService, uc userpb.UserServiceClient) *Handler {
  return &Handler{svc: svc, userClient: uc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
  // 1. Проверить пользователя:
  if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: req.UserId}); err != nil {
    return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
  }
  // 2. Внутренняя логика:
  t, err := h.svc.CreateTask(task.Task{UserID: req.UserId, Title: req.Title})
  if err != nil {
    return nil, err
  }
  // 3. Ответ:
  return &taskpb.CreateTaskResponse{Task: &taskpb.Task{Id: t.ID, UserId: t.UserID, Title: t.Title, IsDone: t.IsDone}}, nil
}

func (h *Handler) GetUser(_ context.Context, req *userpb.User) (*userpb.User, error) {

	findUser, err := h.svc.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &userpb.User{Id: findUser.ID, Email: findUser.Email}, nil
}

func (h *Handler) UpdateUser(_ context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {

	if req.Email == "" {
		return nil, errors.New("email can't be empty")
	}
	id := req.Id
	newUser := user.User{Email: req.Email}

	updatedUser, err := h.svc.UpdateUserByID(id, newUser)
	if err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponse{User: &userpb.User{Id: updatedUser.ID, Email: updatedUser.Email}}, nil
}

func (h *Handler) DeleteUser(_ context.Context, req *userpb.User) (*emptypb.Empty, error) {

	id := req.Id

	err := h.svc.DeleteUserByID(id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) ListUsers(_ context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {

	users, err := h.svc.GetAllUsers()
	if err != nil {
		return nil, err
	}

	arrUsers := make([]*userpb.User, 0, len(users))

	for _, val := range users {
		oneUser := userpb.User{Id: val.ID, Email: val.Email}
		arrUsers = append(arrUsers, &oneUser)
	}

	return &userpb.ListUsersResponse{Users: arrUsers}, nil
}
