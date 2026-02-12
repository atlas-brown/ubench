package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/boutique"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func placeOrder(ctx context.Context, req *boutique.PlaceOrderRequest) *boutique.PlaceOrderResponse {
	order := boutique.PlaceOrder(ctx, req.UserId, req.UserCurrency, req.Address, req.Email, req.CreditCard)
	resp := boutique.PlaceOrderResponse{Order: order}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/place_order", wrappers.Wrapper[boutique.PlaceOrderRequest, boutique.PlaceOrderResponse](placeOrder))
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
