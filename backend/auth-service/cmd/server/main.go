package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/app"
)

func main() {
	a, err := app.NewServerApp()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		os.Exit(1)
	}

	<-ctx.Done()
	a.Stop(context.Background())
}
