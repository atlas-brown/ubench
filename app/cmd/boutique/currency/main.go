package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/boutique"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net"
	"net/http"
	"runtime"
	"os"
	"github.com/goccy/go-json"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func setCurrency(ctx context.Context, req *boutique.SetCurrencySupportRequest) *boutique.SetCurrencySupportResponse {
	ok := boutique.SetCurrencySupport(ctx, req.Currency)
	resp := boutique.SetCurrencySupportResponse{Ok: ok}
	return &resp
}

func getCurrencies(ctx context.Context, req *boutique.GetSupportedCurrenciesRequest) *boutique.GetSupportedCurrenciesResponse {
	currencies := boutique.GetSupportedCurrencies(ctx)
	resp := boutique.GetSupportedCurrenciesResponse{Currencies: currencies}
	return &resp
}

func convertCurrency(ctx context.Context, req *boutique.ConvertCurrencyRequest) *boutique.ConvertCurrencyResponse {
	amount := boutique.ConvertCurrency(ctx, req.Amount, req.ToCurrency)
	resp := boutique.ConvertCurrencyResponse{Amount: amount}
	return &resp
}

func initCurrencies(ctx context.Context, req *boutique.InitCurrencyRequest) *boutique.InitCurrencyResponse {
	boutique.InitCurrencies(ctx, req.Currencies)
	resp := boutique.InitCurrencyResponse{Ok: "OK"}
	return &resp
}

func loadCurrencies(ctx context.Context) []boutique.Currency {
	// List directory
	var currencies []boutique.Currency
	catalogJSON, err := os.ReadFile("/app/cmd/boutique/currency_conversion.json")
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(catalogJSON, &data); err != nil {
		panic(err)
	}
	for key, value := range data {
		currency := boutique.Currency{
			CurrencyCode: key,
			Rate:         fmt.Sprintf("%v", value),
		}
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/set_currency", wrappers.Wrapper[boutique.SetCurrencySupportRequest, boutique.SetCurrencySupportResponse](setCurrency))
	http.HandleFunc("/init_currencies", wrappers.Wrapper[boutique.InitCurrencyRequest, boutique.InitCurrencyResponse](initCurrencies))
	http.HandleFunc("/ro_get_currencies", wrappers.Wrapper[boutique.GetSupportedCurrenciesRequest, boutique.GetSupportedCurrenciesResponse](getCurrencies))
	http.HandleFunc("/ro_convert_currency", wrappers.Wrapper[boutique.ConvertCurrencyRequest, boutique.ConvertCurrencyResponse](convertCurrency))
	boutique.InitAllCurrencies(context.Background(), loadCurrencies(context.Background()))
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
