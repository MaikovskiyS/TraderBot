package strategyservice

import "github.com/MaikovskiyS/TraderBot/internal/trader/providers/indicators"

const (
	StrategySupportRisestence = "SupportRisestence"
)

type Indicators interface {
	Cmo(open []float64, close []float64) []float64
	Rsi(values []float64, period int) []float64
}
type strateger struct {
	indicators Indicators
	*supportResistanceStrategy
	*bidsAsksStrategy
}

func New() *strateger {
	return &strateger{
		indicators: indicators.New(),
		supportResistanceStrategy: &supportResistanceStrategy{
			LookbackPeriod: 20,
			VolLen:         2,
			Interval:       "5",
		},
		bidsAsksStrategy: &bidsAsksStrategy{},
	}
}
