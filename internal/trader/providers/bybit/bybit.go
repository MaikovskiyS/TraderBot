package bybit_provider

import (
	"github.com/hirokisan/bybit/v2"
	"github.com/rs/zerolog"
)

type BybitProvider struct {
	BybitClient *bybit.Client
	Log         zerolog.Logger
}

func New(client *bybit.Client, log zerolog.Logger) *BybitProvider {
	return &BybitProvider{
		BybitClient: client,
		Log:         log,
	}
}
