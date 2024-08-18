package traidingservice

import (
	"testing"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCalculateQty(t *testing.T) {

	testCases := []struct {
		name      string
		balance   float64
		stopLoss  float64
		openPrice float64
		expResult string
	}{
		{
			name:      "buy",
			balance:   100,
			stopLoss:  60000,
			openPrice: 60100,
			expResult: "0.0017",
		},
		{
			name:      "buy",
			balance:   68,
			stopLoss:  59999,
			openPrice: 60456,
			expResult: "0.0011",
		},
		{
			name:      "buy",
			balance:   32,
			stopLoss:  8.43,
			openPrice: 10.00,
			expResult: "0.41",
		},
		{
			name:      "sell",
			balance:   50,
			stopLoss:  0.000863,
			openPrice: 0.000820,
			expResult: "23255",
		},
		{
			name:      "buy",
			balance:   67.9100874,
			stopLoss:  0.2227,
			openPrice: 0.2240,
			expResult: "303",
		},
	}

	for _, test := range testCases {
		tc := test

		t.Run(tc.name, func(t *testing.T) {

			qty := calculateOrderQuantity(tc.balance, tc.openPrice, tc.stopLoss)
			assert.Equal(t, tc.expResult, qty)
		})
	}
}

func TestStopLoss(t *testing.T) {
	testCases := []struct {
		name      string
		candle    *domain.Candle
		side      domain.Side
		precision int
		expResult float64
	}{
		{
			name:      "buy",
			candle:    &domain.Candle{Low: 15, High: 15.5},
			side:      domain.BuySide,
			precision: 2,
			expResult: 14.5,
		},
		{
			name:      "buy",
			candle:    &domain.Candle{Low: 59876, High: 60321},
			side:      domain.BuySide,
			precision: 2,
			expResult: 59431,
		},
		{
			name:      "sell",
			candle:    &domain.Candle{Low: 0.0087, High: 0.0095},
			side:      domain.SellSide,
			precision: 4,
			expResult: 0.010373563218390804,
		},
	}

	for _, test := range testCases {
		tc := test

		t.Run(tc.name, func(t *testing.T) {

			stoploss := calculateOrderStopLoss(tc.candle, tc.side)
			assert.Equal(t, tc.expResult, stoploss)
		})
	}
}
