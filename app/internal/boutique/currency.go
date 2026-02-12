package boutique

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

const (
	debug_currency = false
)

var allCurrencies sync.Map

func InitAllCurrencies(ctx context.Context, currencies []Currency) {
	allCurrencies = sync.Map{}
	for _, currency := range currencies {
		SetCurrencySupport(ctx, currency)
	}
}

func _carry(units float64, nanos float64) carry {
	fractionSize := math.Pow(10, 9)
	nanos += fractionSize
	_units := math.Floor(units) + math.Floor(nanos/fractionSize)
	_nanos := float64(int64(math.Round(nanos)) % int64(fractionSize))
	return carry{Units: _units, Nanos: _nanos}
}

func _MoneyToString(m Money) string {
	nanosStr := strconv.FormatInt(int64(m.Nanos), 10)
	nanosStr = strings.Repeat("0", 9-len(nanosStr)) + nanosStr
	return fmt.Sprintf("%v.%v %v", m.Units, nanosStr, m.Currency)
}

func SetCurrencySupport(ctx context.Context, currency Currency) bool {
	// state.SetState(ctx, currency.CurrencyCode, currency)
	// allCurrencies[currency.CurrencyCode] = currency
	allCurrencies.Store(currency.CurrencyCode, currency)
	return true
}

func InitCurrencies(ctx context.Context, currencies []Currency) {

	for _, currency := range currencies {
		SetCurrencySupport(ctx, currency)
	}
}

func ConvertCurrency(ctx context.Context, amount Money, toCurrency string) Money {
	fromRate, ok := allCurrencies.Load(amount.Currency)
	if !ok {
		panic(fmt.Errorf("currency %s not found", amount.Currency))
	}
	fromRate_, _ := fromRate.(Currency)
	fromRate64, err := strconv.ParseFloat(fromRate_.Rate, 64)
	if err != nil {
		panic(err)
	}

	// Convert: from_currency --> EUR
	euros := _carry(float64(amount.Units)/fromRate64, float64(amount.Nanos)/fromRate64)

	euros.Nanos = math.Round(euros.Nanos)

	// Convert: EUR --> to_currency
	// toRate, err := state.GetState[Currency](ctx, toCurrency)
	// toRate, ok := allCurrencies[toCurrency]
	toRate, ok := allCurrencies.Load(toCurrency)
	toRate_, _ := toRate.(Currency)
	if !ok {
		panic(fmt.Errorf("currency %s not found", toCurrency))
	}

	toRate64, err := strconv.ParseFloat(toRate_.Rate, 64)
	if err != nil {
		panic(err)
	}
	_result := _carry(euros.Units*toRate64, euros.Nanos*toRate64)

	_result.Units = math.Floor(_result.Units)
	_result.Nanos = math.Floor(_result.Nanos)
	result := Money{Currency: toCurrency, Units: int32(_result.Units), Nanos: int64(_result.Nanos)}

	return result
}

func GetSupportedCurrencies(ctx context.Context) []Currency {
	if debug_currency { fmt.Println("GetSupportedCurrencies ") }

	var currencies []Currency
	allCurrencies.Range(func(key, value interface{}) bool {
		currencies = append(currencies, value.(Currency))
		return true
	})
	return currencies
}
