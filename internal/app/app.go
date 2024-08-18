package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MaikovskiyS/TraderBot/internal/config"
	bybit_provider "github.com/MaikovskiyS/TraderBot/internal/trader/providers/bybit"
	strategyservice "github.com/MaikovskiyS/TraderBot/internal/trader/services/strategy_service"
	traidingservice "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/MaikovskiyS/TraderBot/internal/trader/usecases"
	"github.com/MaikovskiyS/TraderBot/pkg/logger"
	"github.com/hirokisan/bybit/v2"
	"github.com/rs/zerolog"
)

func Run() error {
	cfg, err := config.LoadBybitConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := logger.New(zerolog.DebugLevel)
	logger.Info().Msg("init service")

	bybitCl := bybit.NewClient().
		WithAuth(cfg.PublicKey, cfg.PrivateKey).
		WithBaseURL(cfg.BaseURL)

	bybitProvider := bybit_provider.New(
		bybitCl,
		logger)

	trading := traidingservice.New(bybitProvider, logger)
	strategy := strategyservice.New()
	usecases := usecases.New(trading, strategy, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("shutdown signal received, stopping worker")
				return
			case <-ticker.C:

				bybitCl.SyncServerTime()

				logger.Debug().Msg("manage positions")
				err = usecases.ManagePositions(ctx)
				if err != nil {
					if errors.Is(err, traidingservice.ErrNoPositions) {
						logger.Debug().Msg("run strategy")
						err = usecases.RunSupportResistance(ctx)
						if err != nil {
							logger.Error().Err(err).Send()
						}
					} else {
						logger.Error().Err(err).Send()
					}
				}
			}
		}
	}()

	<-sigCh
	logger.Info().Msg("shutting down gracefully")

	cancel()
	wg.Wait()

	return nil
}

// складывать асинхронно собирать значения свечей для индикаторов, расчеты индикаторов складывать в мапу.
// Написать глиента к мапе. для получения различных значений индикаторов
