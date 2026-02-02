package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ailinykh/reposter/v3/pkg/telegram"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	logger := NewLogger()
	bot, err := telegram.NewBot(ctx,
		telegram.WithToken(config.BotToken),
		telegram.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}
	logger.Info("telegram bot initialized", "username", bot.Username)

	lightsail := NewLightsailAPI(config.KeyID, config.KeySecret)
	instances, err := lightsail.GetServers(ctx)
	if err != nil {
		panic(err)
	}
	logger.Info("lightsail api initialized", "instance_count", len(instances))

	var lastReboot = time.Now()
	var isUp = false
	for {
		select {
		case <-ctx.Done():
			logger.Info("attempt to shutdown gracefully...")
			return

		case <-time.After(15 * time.Second):
			e := healthcheck(ctx, config.URL, logger)
			if e != nil {
				logger.Error("healthcheck failed", "error", e)
				if isUp {
					instanceName := *instances[0].Name
					logger.Info("performing reboot", "instance_name", instanceName)
					_, err := lightsail.RebootServer(ctx, instances[0])
					if err != nil {
						logger.Error("failed to reboot instance", "error", err)
					} else {
						if _, err := bot.SendMessage(ctx, &telegram.SendMessageParams{
							ChatID:    config.ChannelID,
							Text:      fmt.Sprintf("🌈 Instance <b>%s</b> rebooted\n\nuptime is <b>%v</b>", instanceName, time.Since(lastReboot.Round(time.Second))),
							ParseMode: telegram.ParseModeHTML,
						}); err != nil {
							logger.Error("failed to send message", "error", err)
						}
						logger.Info("✅ instance rebooted!", "uptime", time.Since(lastReboot))
						lastReboot = time.Now()
					}
				}
			}
			isUp = e == nil
		}
	}
}

func healthcheck(ctx context.Context, url string, logger *slog.Logger) error {
	logger.Info("healthcheck", "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	logger.Info("it's ok!", "status", res.Status)
	return nil
}
