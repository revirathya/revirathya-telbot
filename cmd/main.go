package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// Internal
	"github.com/revirathya/revirathya-telbot/cmd/api"
	"github.com/revirathya/revirathya-telbot/internal/config"
	"github.com/revirathya/revirathya-telbot/internal/delivery"
	botDlv "github.com/revirathya/revirathya-telbot/internal/delivery/bot"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/clients"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/repository"
	"github.com/revirathya/revirathya-telbot/internal/usecase"
)

// version is set at build time: go build -ldflags "-X main.version=0.2.0"
var version = "dev"

func main() {
	// Get Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	// Run
	if err := run(cfg); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	// Set logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Define context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Get Telegram Bot
	tbh, err := telegram.GetBotWebhook(cfg.TelegramBotToken, cfg.TelegramWebhookURL)
	if err != nil {
		return fmt.Errorf("failed to get webhook: %w", err)
	}

	// Set gRPC Clients
	finance, err := clients.NewFinanceClient(cfg.FinanceHost, cfg.FinanceApiKey)
	if err != nil {
		return fmt.Errorf("failed to create finance grpc client: %w", err)
	}
	defer finance.Close()

	// Set repo and usecase
	repo := repository.BootstrapRepository(finance)
	uc := usecase.BootstrapUsecase(repo)

	// Setup server
	httpSrvLayer := delivery.Setup(cfg.AppEnv, cfg.AppPort, version)

	// Configure bot
	if err := botDlv.RegisterHandler(ctx, httpSrvLayer, tbh, uc); err != nil {
		return fmt.Errorf("failed to register bot delivery handle: %w", err)
	}

	g := api.RunServers(ctx, httpSrvLayer)
	slog.Info(
		"server started",
		"version", version,
		"http_port", cfg.AppPort,
	)
	if err := g.Wait(); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
