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
	createdTask, err := h.svc.CreateTask(task.Task{UserID: req.UserId, Title: req.Title})
	if err != nil {
		return nil, err
	}
	// 3. Ответ:
	return &taskpb.CreateTaskResponse{Task: &taskpb.Task{Id: createdTask.ID, UserId: createdTask.UserID, Title: createdTask.Title}}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.Task) (*taskpb.Task, error) {

	findTask, err := h.svc.GetTaskByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &taskpb.Task{Id: findTask.ID, Title: findTask.Title, UserId: findTask.UserID}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {

	if req.Title == "" {
		return nil, errors.New("title can't be empty")
	}
	if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: req.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	id := req.Id
	newTask := task.Task{Title: req.Title, UserID: req.UserId}

	updatedTask, err := h.svc.UpdateTaskByID(id, newTask)
	if err != nil {
		return nil, err
	}

	return &taskpb.UpdateTaskResponse{Task: &taskpb.Task{Id: updatedTask.ID, Title: updatedTask.Title, UserId: updatedTask.UserID}}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.Task) (*emptypb.Empty, error) {

	id := req.Id

	err := h.svc.DeleteTaskByID(id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {

	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, err
	}

	arrTasks := make([]*taskpb.Task, 0, len(tasks))

	for _, val := range tasks {
		oneTask := taskpb.Task{Id: val.ID, Title: val.Title, UserId: val.UserID}
		arrTasks = append(arrTasks, &oneTask)
	}

	return &taskpb.ListTasksResponse{Tasks: arrTasks}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksByUserResponse, error) {

	if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: req.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	allUserTasks, err := h.svc.GetTasksByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	userTasks := make([]*taskpb.Task, 0, len(allUserTasks))
	
	for _, val := range allUserTasks {
		oneUserTask := taskpb.Task{Id: val.ID, Title: val.Title, UserId: val.UserID}
		userTasks = append(userTasks, &oneUserTask)
	}

	return &taskpb.ListTasksByUserResponse{Tasks: userTasks}, nil
}
