package main

import (
	"context"
	"fmt"
	// "github.com/atlas/slowpoke/pkg/cm"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net/http"
	"runtime"
	"github.com/atlas/slowpoke/internal/social"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func readUserTimeline(ctx context.Context, req *social.ReadUserTimelineRequest) *social.ReadUserTimelineResponse {
	posts := social.ReadUserTimeline(ctx, req.UserId)
	//fmt.Printf("Posts read: %+v\n", posts)
	resp := social.ReadUserTimelineResponse{Posts: posts}
	return &resp
}

func writeUserTimeline(ctx context.Context, req *social.WriteUserTimelineRequest) *string {
	social.WriteUserTimeline(ctx, req.UserId, req.PostIds)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read_user_timeline", wrappers.Wrapper[social.ReadUserTimelineRequest, social.ReadUserTimelineResponse](readUserTimeline))
	http.HandleFunc("/write_user_timeline", wrappers.Wrapper[social.WriteUserTimelineRequest, string](writeUserTimeline))
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
