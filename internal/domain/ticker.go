package domain

import (
	"sort"
	"strconv"
)

type Ticker struct {
	Symbol            string
	LastPrice         string
	IndexPrice        string
	MarkPrice         string
	PrevPrice24h      string
	Price24hPcnt      string
	HighPrice24h      string
	LowPrice24h       string
	PrevPrice1h       string
	OpenInterest      string
	OpenInterestValue string
	Turnover24h       string
	Volume24h         string
	Ask1Size          string
	Bid1Price         string
	Ask1Price         string
	Bid1Size          string
	Ask1Iv            string
	Bid1Iv            string
	MarkIv            string
	Precision         int
}

type TickerPrecision struct {
	Symbol    string
	Precision int
}

func (ti Ticker) LastPriceFloat64() float64 {
	lastPrice, _ := strconv.ParseFloat(ti.LastPrice, 64)
	return lastPrice
}

type Tickers []*Ticker

func (t Tickers) Get10MostVolume24h() Tickers {
	// Сортируем тикеры по всем трём критериям (Volume24h, Price24hPcnt, OpenInterest)
	sort.Slice(t, func(i, j int) bool {
		volumeI, _ := strconv.ParseFloat(t[i].Volume24h, 64)
		volumeJ, _ := strconv.ParseFloat(t[j].Volume24h, 64)
		priceI, _ := strconv.ParseFloat(t[i].LastPrice, 64)
		priceJ, _ := strconv.ParseFloat(t[j].LastPrice, 64)
		if volumeI != volumeJ {
			return volumeI*priceI > volumeJ*priceJ
		}

		priceChangeI, _ := strconv.ParseFloat(t[i].Price24hPcnt, 64)
		priceChangeJ, _ := strconv.ParseFloat(t[j].Price24hPcnt, 64)
		if priceChangeI != priceChangeJ {
			return priceChangeI > priceChangeJ
		}

		openInterestI, _ := strconv.ParseFloat(t[i].OpenInterest, 64)
		openInterestJ, _ := strconv.ParseFloat(t[j].OpenInterest, 64)
		return openInterestI > openInterestJ
	})

	// Возвращаем топ 40 монет
	if len(t) > 20 {
		return t[:40]
	}
	return t
}
