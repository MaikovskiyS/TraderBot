package indicators

import "github.com/MaikovskiyS/TraderBot/internal/domain"

type MovingAverageType string

const (
	SMA MovingAverageType = "SMA"
	EMA MovingAverageType = "EMA"
)

func calculateMA(candles []*domain.Candle, length int, maType MovingAverageType) []float64 {
	var result []float64
	switch maType {
	case SMA:
		for i := length - 1; i < len(candles); i++ {
			sum := 0.0
			for j := i - length + 1; j <= i; j++ {
				sum += candles[j].Close
			}
			result = append(result, sum/float64(length))
		}
	case EMA:
		// EMA calculation (simplified version)
		ema := make([]float64, len(candles))
		k := 2.0 / float64(length+1)
		ema[length-1] = candles[length-1].Close
		for i := length; i < len(candles); i++ {
			ema[i] = (candles[i].Close-k)*ema[i-1] + k*candles[i].Close
		}
		result = ema
	}
	return result
}
