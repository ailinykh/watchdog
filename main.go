package main

import (
	"context"
	"fmt"
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

	instance := instances[0]
	var rebootStartedAt = time.Now()
	var rebootFinishedAt = time.Now()
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
				if !isUp {
					break
				}
				// looks like server went down
				logger.Info("perform reboot", "instance_name", *instance.Name, "error", e)
				if _, err := lightsail.RebootServer(ctx, instance); err != nil {
					logger.Error("failed to reboot instance", "error", err)
					break
				}

				if _, err = bot.SendMessage(ctx, &telegram.SendMessageParams{
					ChatID:    config.ChannelID,
					Text:      fmt.Sprintf("🔴 <b>%s</b> is offline\nInstance last uptime is: <b>%v</b>", *instance.Name, time.Since(rebootFinishedAt).Round(time.Second)),
					ParseMode: telegram.ParseModeHTML,
				}); err != nil {
					logger.Error("failed to send message", "error", err)
				}

				logger.Info("🌈 waiting for instance to recover", "last_uptime", time.Since(rebootFinishedAt))
				rebootStartedAt = time.Now()
			} else if !isUp && rebootStartedAt.Round(time.Millisecond) != rebootFinishedAt.Round(time.Millisecond) {
				// server recovered
				if _, err = bot.SendMessage(ctx, &telegram.SendMessageParams{
					ChatID: config.ChannelID,
					Text: fmt.Sprintf(
						"🟢 <b>%s</b> is online\nInstance recovered after downtime: <b>%v</b>",
						*instance.Name,
						time.Since(rebootStartedAt).Round(time.Second),
					),
					ParseMode: telegram.ParseModeHTML,
				}); err != nil {
					logger.Error("failed to send message", "error", err)
				}

				logger.Info("🚀 instance rebooted!", "downtime", time.Since(rebootStartedAt))
				rebootFinishedAt = time.Now()
			}

			isUp = e == nil
		}
	}
}
