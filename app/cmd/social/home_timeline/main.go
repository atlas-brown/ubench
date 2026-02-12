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

func readHomeTimeline(ctx context.Context, req *social.ReadHomeTimelineRequest) *social.ReadHomeTimelineResponse {
	posts := social.ReadHomeTimeline(ctx, req.UserId)
	//fmt.Printf("Posts read: %+v\n", posts)
	resp := social.ReadHomeTimelineResponse{Posts: posts}
	return &resp
}

func writeHomeTimeline(ctx context.Context, req *social.WriteHomeTimelineRequest) *string {
	social.WriteHomeTimeline(ctx, req.UserId, req.PostIds)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read_home_timeline", wrappers.Wrapper[social.ReadHomeTimelineRequest, social.ReadHomeTimelineResponse](readHomeTimeline))
	http.HandleFunc("/write_home_timeline", wrappers.Wrapper[social.WriteHomeTimelineRequest, string](writeHomeTimeline))
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
