package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"next-video-golang/internal/account"
	"next-video-golang/internal/config"
	"next-video-golang/internal/database"
	"next-video-golang/internal/httpapi"
	"next-video-golang/internal/interaction"
	"next-video-golang/internal/video"
)

// main 启动 API 服务；读取本地配置、连接数据库并监听 HTTP 请求。
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	videoRepository := video.NewPostgresRepository(pool)
	videoService := video.NewService(videoRepository)
	accountRepository := account.NewPostgresRepository(pool)
	accountService := account.NewService(accountRepository, time.Now)
	interactionRepository := interaction.NewPostgresRepository(pool)
	interactionService := interaction.NewService(interactionRepository, time.Now, nil)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewServerWithServices(videoService, accountService, interactionService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("next-video api listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}
