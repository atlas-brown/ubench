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

func sendEmail(ctx context.Context, req *boutique.SendOrderConfirmationRequest) *boutique.SendOrderConfirmationResponse {
	ok := boutique.SendConfirmation(ctx, req.Email, req.Order)
	resp := boutique.SendOrderConfirmationResponse{Ok: ok}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_send_email", wrappers.Wrapper[boutique.SendOrderConfirmationRequest, boutique.SendOrderConfirmationResponse](sendEmail))
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	panic(http.Serve(listener, nil))
}
