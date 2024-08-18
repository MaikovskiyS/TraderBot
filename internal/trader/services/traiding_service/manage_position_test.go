package traidingservice

// func TestCalcMovingPercent(t *testing.T) {
// 	testCases := []struct {
// 		name      string
// 		position  *domain.Position
// 		expResult float64
// 	}{
// 		{
// 			name:      "buy",
// 			position:  &domain.Position{AvgPrice: 231.00, StopLoss: 3333.00},
// 			expResult: 14.5,
// 		},
// 		{
// 			name:     "buy",
// 			position: &domain.Position{AvgPrice: 231.00, StopLoss: 3333.00},

// 			expResult: 59431,
// 		},
// 		{
// 			name:      "sell",
// 			position:  &domain.Position{AvgPrice: 231.00, StopLoss: 3333.00},
// 			expResult: 0.010373563218390804,
// 		},
// 	}

// 	for _, test := range testCases {
// 		tc := test

// 		t.Run(tc.name, func(t *testing.T) {

// 			stoploss := calculateOrderStopLoss(tc.candle, tc.side)
// 			assert.Equal(t, tc.expResult, stoploss)
// 		})
// 	}
// }
