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

func addItemToCart(ctx context.Context, req *boutique.AddItemRequest) *boutique.AddItemResponse {
	ok := boutique.AddItem(ctx, req.UserId, req.ProductId, req.Quantity)
	resp := boutique.AddItemResponse{Ok: ok}
	return &resp
}

func getCart(ctx context.Context, req *boutique.GetCartRequest) *boutique.GetCartResponse {
	cart := boutique.GetCart(ctx, req.UserId)
	resp := boutique.GetCartResponse{Cart: cart}
	return &resp
}

func emptyCart(ctx context.Context, req *boutique.EmptyCartRequest) *boutique.EmptyCartResponse {
	ok := boutique.EmptyCart(ctx, req.UserId)
	resp := boutique.EmptyCartResponse{Ok: ok}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/add_item", wrappers.Wrapper[boutique.AddItemRequest, boutique.AddItemResponse](addItemToCart))
	http.HandleFunc("/ro_get_cart", wrappers.Wrapper[boutique.GetCartRequest, boutique.GetCartResponse](getCart))
	http.HandleFunc("/empty_cart", wrappers.Wrapper[boutique.EmptyCartRequest, boutique.EmptyCartResponse](emptyCart))
	boutique.CartInit()
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
