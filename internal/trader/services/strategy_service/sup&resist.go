package strategyservice

import (
	"math"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

type supportResistanceBox struct {
	Ticker       string
	OpenPrice    float64
	TargetCandle *domain.Candle
	Side         domain.Side
}

type supportResistanceStrategy struct {
	LookbackPeriod int
	VolLen         int
	Interval       string
}

func (s *supportResistanceStrategy) ApplySupportResistance(data []*domain.TickerCandels) (*StrategyResponse, error) {
	var trades []*Trade

	for _, ticker := range data {
		buyBox, sellBox := calculateSupportResistance(ticker.Candels)
		lastCandle := ticker.Candels[0]

		for _, box := range []*supportResistanceBox{buyBox, sellBox} {
			if box.OpenPrice == 0 {
				continue
			}

			var closePrice float64
			var profitPercent float64

			// Определение ClosePrice и проверка изменения цены на основе Side
			var priceDifferencePercent float64
			switch box.Side {
			case domain.BuySide:
				if lastCandle.Low <= box.OpenPrice {
					continue
				}
				closePrice = findPreviousHigh(ticker.Candels, box.TargetCandle)
				priceDifferencePercent = ((lastCandle.Open - box.TargetCandle.Low) / box.TargetCandle.Low) * 100
			case domain.SellSide:
				if lastCandle.High >= box.OpenPrice {
					continue
				}
				closePrice = findPreviousLow(ticker.Candels, box.TargetCandle)
				priceDifferencePercent = ((lastCandle.Open - box.TargetCandle.High) / box.TargetCandle.High) * 100
			}

			// Проверка разницы в процентах
			if math.Abs(priceDifferencePercent) > 1 {
				continue
			}

			// Расчет прибыли в процентах
			if box.OpenPrice != 0 && closePrice != 0 {
				profitPercent = math.Abs(((closePrice - box.OpenPrice) / box.OpenPrice) * 100)
			}

			trades = append(trades, &Trade{
				Symbol:        ticker.Symbol,
				Side:          box.Side,
				OpenPrice:     box.OpenPrice,
				ClosePrice:    &closePrice,
				TargetCandle:  box.TargetCandle,
				ProfitPercent: profitPercent,
			})
		}
	}

	// Найти лучший вариант на основе profitPercent
	if len(trades) == 0 {
		return nil, nil
	}

	bestTrade := trades[0]
	for _, trade := range trades[1:] {
		if trade.ProfitPercent > bestTrade.ProfitPercent {
			bestTrade = trade
		}
	}

	return &StrategyResponse{Trade: bestTrade}, nil
}

func calculateSupportResistance(candles []*domain.Candle) (buyBox *supportResistanceBox, sellBox *supportResistanceBox) {
	var lowestPivotLow float64 = math.Inf(1)
	var highestPivotHigh float64 = math.Inf(-1)
	var lowestPivotCandle, highestPivotCandle *domain.Candle

	for i := 1; i < len(candles)-1; i++ {
		// Определяем пивоты
		pivotHigh := getPivotHigh(candles, i)
		pivotLow := getPivotLow(candles, i)

		// Обновляем самый низкий пивот для покупки
		if pivotLow != 0 && pivotLow < lowestPivotLow {
			lowestPivotLow = pivotLow
			lowestPivotCandle = candles[i]
		}

		// Обновляем самый высокий пивот для продажи
		if pivotHigh != 0 && pivotHigh > highestPivotHigh {
			highestPivotHigh = pivotHigh
			highestPivotCandle = candles[i]
		}
	}

	// Если найден самый низкий пивот, создаем структуру для покупки
	if lowestPivotLow != math.Inf(1) {
		buyBox = &supportResistanceBox{
			OpenPrice:    lowestPivotCandle.GetMiddlePrice(),
			TargetCandle: lowestPivotCandle,
			Side:         domain.BuySide,
		}
	}

	// Если найден самый высокий пивот, создаем структуру для продажи
	if highestPivotHigh != math.Inf(-1) {
		sellBox = &supportResistanceBox{
			OpenPrice:    highestPivotCandle.GetMiddlePrice(),
			TargetCandle: highestPivotCandle,
			Side:         domain.SellSide,
		}
	}

	return buyBox, sellBox
}

func findPreviousLow(candles []*domain.Candle, targetCandle *domain.Candle) float64 {
	var previousLow float64

	// Проходим по свечам до целевой свечи
	for i, candle := range candles {
		// Если дошли до целевой свечи, прекращаем цикл
		if candle.Time == targetCandle.Time {
			break
		}

		// Инициализируем previousLow первым найденным значением Low
		if i == 0 {
			previousLow = candle.Low
		}

		// Находим минимум до целевой свечи
		if candle.Low < previousLow {
			previousLow = candle.Low
		}
	}

	return previousLow
}

func findPreviousHigh(candles []*domain.Candle, targetCandle *domain.Candle) float64 {
	var previousHigh float64

	// Проходим по свечам до целевой свечи
	for _, candle := range candles {
		// Если дошли до целевой свечи, прекращаем цикл
		if candle.Time == targetCandle.Time {
			break
		}

		// Находим максимум до целевой свечи
		if candle.High > previousHigh {
			previousHigh = candle.High
		}
	}

	return previousHigh
}

func getPivotHigh(candles []*domain.Candle, index int) float64 {
	if index == 0 || index == len(candles)-1 {
		return 0
	}
	if candles[index].High > candles[index-1].High && candles[index].High > candles[index+1].High {
		return candles[index].High
	}
	return 0
}

func getPivotLow(candles []*domain.Candle, index int) float64 {
	if index == 0 || index == len(candles)-1 {
		return 0
	}
	if candles[index].Low < candles[index-1].Low && candles[index].Low < candles[index+1].Low {
		return candles[index].Low
	}
	return 0
}
