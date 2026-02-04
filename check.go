package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

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
