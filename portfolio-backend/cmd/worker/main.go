package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kojima1128/portfolio-backend/internal/infrastructure"
	mysqlrepo "github.com/kojima1128/portfolio-backend/internal/infrastructure/mysql"
	"github.com/kojima1128/portfolio-backend/internal/service"
)

func main() {
	db := infrastructure.NewDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	userRepo := mysqlrepo.NewUserRepository(db)
	// TODO: Use userService in database monitoring logic
	userService := service.NewUserService(userRepo)
	_ = userService

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		log.Println("monitoring goroutine started")
		for {
			select {
			case <-ctx.Done():
				log.Println("monitoring goroutine stopped")
				return
			case <-ticker.C:
				// TODO: Implement database monitoring logic
			}
		}
	}()

	log.Println("worker started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down worker...")
	cancel()
	wg.Wait()
	log.Println("worker stopped")
}
