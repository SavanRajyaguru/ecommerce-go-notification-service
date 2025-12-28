package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/database"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/event"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/notification"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/processor"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/pkg/utils"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/repository"
)

func main() {
	// 1. Initialize Logger
	utils.InitLogger()
	log.Println("Starting Notification Service Worker...")

	// 2. Load Config
	config.LoadConfig()

	// 3. Connect to MongoDB
	database.ConnectDB()
	defer database.DisconnectDB()

	// 4. Initialize Components
	notifRepo := repository.NewNotificationRepository()
	emailSender := notification.NewEmailSender()
	smsSender := notification.NewSMSSender()

	proc := processor.NewNotificationProcessor(notifRepo, emailSender, smsSender)
	consumer := event.NewConsumer(proc)

	// 5. Start Consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Start(ctx)

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	cancel() // Stop consumer
	log.Println("Worker stopped")
}
