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

func getQuote(ctx context.Context, req *boutique.GetQuoteRequest) *boutique.GetQuoteResponse {
	// slowpoke.SlowpokeCheck("getQuote")
	quote := boutique.GetQuote(ctx, req.Items)
	resp := boutique.GetQuoteResponse{CostUsd: quote}
	return &resp
}

func shipOrder(ctx context.Context, req *boutique.ShipOrderRequest) *boutique.ShipOrderResponse {
	// slowpoke.SlowpokeCheck("shipOrder")
	id := boutique.ShipOrder(ctx, req.Address, req.Items)
	resp := boutique.ShipOrderResponse{TrackingId: id}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_get_quote", wrappers.Wrapper[boutique.GetQuoteRequest, boutique.GetQuoteResponse](getQuote))
	http.HandleFunc("/ship_order", wrappers.Wrapper[boutique.ShipOrderRequest, boutique.ShipOrderResponse](shipOrder))
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
