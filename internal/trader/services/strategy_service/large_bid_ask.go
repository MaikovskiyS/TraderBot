package strategyservice

import (
	"context"
	"fmt"
	"math"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

type bidsAsksStrategy struct {
}

type Level struct {
	Price    float64
	Quantity float64
	Side     domain.Side // "Bid" or "Ask"
}

func (s *bidsAsksStrategy) ApplyBidsAsksStrategy(
	ctx context.Context, tickers []*domain.TickerOrderBook) (*StrategyResponse, error) {

	for _, ticker := range tickers {
		treshold := calculateThreshold(ticker.OrderBook)

		levels := findLargeOrders(ticker.OrderBook, treshold)
		mapa := aggregateLevels(ticker.OrderBook)
		fmt.Printf("mapa: %v\n", mapa)

		for _, level := range levels {
			fmt.Printf("level: %v\n", level)
		}
	}

	return nil, nil
}

// aggregateLevels суммирует quantity на округленных уровнях цен
func aggregateLevels(orderBook *domain.OrderBook) map[float64]float64 {
	levelSums := make(map[float64]float64)

	// Проходим по Bids и Asks
	for _, entry := range append(orderBook.Bids, orderBook.Asks...) {
		roundedPrice := customRound(entry.Price)
		levelSums[roundedPrice] += entry.Quantity
	}

	return levelSums
}

func customRound(value float64) float64 {
	// Округляем до целого, если значение больше 1
	if value >= 1 {
		return math.Floor(value)
	}

	// Если значение меньше 1, то сокращаем до определенного количества знаков после запятой
	if value < 1 {
		factor := 1.0

		// Определяем количество знаков после запятой
		for value*factor < 1 {
			factor *= 10
		}

		// Округляем до двух значимых цифр после первой ненулевой
		value = math.Floor(value*factor*10+0.5) / (factor * 10)
	}

	return value
}

func findLargeOrders(orderBook *domain.OrderBook, threshold float64) []Level {
	var levels []Level

	// Analyze Bids
	for _, bid := range orderBook.Bids {
		if bid.Quantity >= threshold {
			levels = append(levels, Level{
				Price:    bid.Price,
				Quantity: bid.Quantity,
				Side:     domain.SellSide,
			})
		}
	}

	// Analyze Asks
	for _, ask := range orderBook.Asks {
		if ask.Quantity >= threshold {
			levels = append(levels, Level{
				Price:    ask.Price,
				Quantity: ask.Quantity,
				Side:     domain.BuySide,
			})
		}
	}

	return levels
}

func calculateThreshold(orderBook *domain.OrderBook) float64 {
	const multiplier = 3
	/*
			Как выбрать multiplier:
		Меньшие значения (например, 1.5): Это более чувствительный подход, который будет учитывать больше заявок как значимые.
		Большие значения (например, 2 или 3): Это более строгий подход, который будет учитывать только самые крупные заявки.
	*/
	totalVolume := 0.0
	totalOrders := 0

	// Summing up the volumes in bids and asks
	for _, bid := range orderBook.Bids {
		totalVolume += bid.Quantity
		totalOrders++
	}
	for _, ask := range orderBook.Asks {
		totalVolume += ask.Quantity
		totalOrders++
	}

	// Calculate the average volume
	averageVolume := totalVolume / float64(totalOrders)

	// Apply the multiplier to determine the threshold
	threshold := averageVolume * multiplier

	return threshold
}
