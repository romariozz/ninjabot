package model

import (
	"fmt"
	"strconv"
	"time"
)

type TelegramSettings struct {
	Enabled bool
	Token   string
	Users   []int
}

type Settings struct {
	Pairs    []string
	Telegram TelegramSettings
}

type Balance struct {
	Tick string
	Free float64
	Lock float64
}

type AssetInfo struct {
	BaseAsset  string
	QuoteAsset string

	MinPrice    float64
	MaxPrice    float64
	MinQuantity float64
	MaxQuantity float64
	StepSize    float64
	TickSize    float64

	// Number of decimal places
	QtyDecimalPrecision   int64
	PriceDecimalPrecision int64
}

type Dataframe struct {
	Pair string

	Close  Series
	Open   Series
	High   Series
	Low    Series
	Volume Series

	Time       []time.Time
	LastUpdate time.Time

	// Custom user metadata
	Metadata map[string]Series
}

type Candle struct {
	Pair      string
	Time      time.Time
	UpdatedAt time.Time
	Open      float64
	Close     float64
	Low       float64
	High      float64
	Volume    float64
	Trades    int64
	Complete  bool
}

func (c Candle) ToSlice(precision int) []string {
	return []string{
		fmt.Sprintf("%d", c.Time.Unix()),
		strconv.FormatFloat(c.Open, 'f', precision, 64),
		strconv.FormatFloat(c.Close, 'f', precision, 64),
		strconv.FormatFloat(c.Low, 'f', precision, 64),
		strconv.FormatFloat(c.High, 'f', precision, 64),
		fmt.Sprintf("%.1f", c.Volume),
		fmt.Sprintf("%d", c.Trades),
	}
}

func (c Candle) Less(j Item) bool {
	diff := j.(Candle).Time.Sub(c.Time)
	if diff < 0 {
		return false
	}
	if diff > 0 {
		return true
	}

	diff = j.(Candle).UpdatedAt.Sub(c.UpdatedAt)
	if diff < 0 {
		return false
	}
	if diff > 0 {
		return true
	}

	return c.Pair < j.(Candle).Pair
}

type Account struct {
	Balances []Balance
}

func (a Account) Balance(tick string) Balance {
	for _, balance := range a.Balances {
		if balance.Tick == tick {
			return balance
		}
	}
	return Balance{}
}

func (a Account) Equity() float64 {
	var total float64

	for _, balance := range a.Balances {
		total += balance.Free
		total += balance.Lock
	}

	return total
}
