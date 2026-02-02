package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	BotToken  string
	ChannelID int64
	KeyID     string
	KeySecret string
	URL       string
}

func NewConfig() (*Config, error) {
	var (
		botToken  string
		channelID int64
		keyID     string
		keySecret string
		url       string
		ok        bool
		err       error
	)

	if botToken, ok = os.LookupEnv("TELEGRAM_BOT_TOKEN"); !ok {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	if channelID, err = strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64); err != nil {
		return nil, fmt.Errorf("failed to parse CHANNEL_ID: %w", err)
	}

	if keyID, ok = os.LookupEnv("KEY_ID"); !ok {
		return nil, fmt.Errorf("KEY_ID not set")
	}

	if keySecret, ok = os.LookupEnv("KEY_SECRET"); !ok {
		return nil, fmt.Errorf("KEY_SECRET not set")
	}

	if url, ok = os.LookupEnv("URL"); !ok {
		return nil, fmt.Errorf("URL not set")
	}

	return &Config{
		BotToken:  botToken,
		ChannelID: channelID,
		KeyID:     keyID,
		KeySecret: keySecret,
		URL:       url,
	}, nil
}
