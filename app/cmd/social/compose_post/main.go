package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/social"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func ComposePost(ctx context.Context, req *social.ComposePostRequest) *string {
	social.ComposePost(ctx, req.Text, req.CreatorId)
	resp := "OK"
	return &resp
}

func ComposePostMulti(ctx context.Context, req *social.ComposePostMultiRequest) *string {
	social.ComposeMulti(ctx, req.Text, req.Number, req.CreatorId)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/compose_post", wrappers.Wrapper[social.ComposePostRequest, string](ComposePost))
	http.HandleFunc("/compose_post_multi", wrappers.Wrapper[social.ComposePostMultiRequest, string](ComposePostMulti))
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
