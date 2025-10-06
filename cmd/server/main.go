package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guarref/tasks-service/internal/database"
	"github.com/guarref/tasks-service/internal/task"
	"github.com/guarref/tasks-service/internal/transport/grpc"
)

func main() {
	
	database.InitDB()
	if err := database.DB.AutoMigrate(&task.Task{}); err != nil {
		log.Fatalf("failed to automigrate tasks table: %v", err)
	}

	// 2. Репозиторий и сервис задач
	repo := task.NewTaskRepository(database.DB)
	svc := task.NewTaskService(repo)

	// 3. Клиент к Users-сервису
	userClient, conn, err := grpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to users: %v", err)
	}
	defer conn.Close()

	// 4. Запуск gRPC Tasks-сервера
	if err := grpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}

	// 5. Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down server...")
}