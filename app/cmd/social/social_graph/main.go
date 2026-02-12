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

func InsertUser(ctx context.Context, req *social.InsertUserRequest) *string {
	social.InsertUser(ctx, req.UserId)
	resp := "OK"
	return &resp
}

func GetFollowers(ctx context.Context, req *social.GetFollowersRequest) *social.GetFollowersResponse {
	followers := social.GetFollowers(ctx, req.UserId)
	resp := social.GetFollowersResponse{
		Followers: followers,
	}
	return &resp
}

func GetFollowees(ctx context.Context, req *social.GetFolloweesRequest) *social.GetFolloweesResponse {
	followees := social.GetFollowees(ctx, req.UserId)
	resp := social.GetFolloweesResponse{
		Followees: followees,
	}
	return &resp
}

func Follow(ctx context.Context, req *social.FollowRequest) *string {
	social.Follow(ctx, req.FollowerId, req.FolloweeId)
	resp := "OK"
	return &resp
}

func FollowMulti(ctx context.Context, req *social.FollowManyRequest) *string {
	social.FollowMulti(ctx, req.UserId, req.FollowerIds, req.FolloweeIds)
	resp := "OK"
	return &resp
}

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/insert_user", wrappers.Wrapper[social.InsertUserRequest, string](InsertUser))
	http.HandleFunc("/ro_get_followers", wrappers.Wrapper[social.GetFollowersRequest, social.GetFollowersResponse](GetFollowers))
	http.HandleFunc("/ro_get_followees", wrappers.Wrapper[social.GetFolloweesRequest, social.GetFolloweesResponse](GetFollowees))
	http.HandleFunc("/follow", wrappers.Wrapper[social.FollowRequest, string](Follow))
	http.HandleFunc("/follow_multi", wrappers.Wrapper[social.FollowManyRequest, string](FollowMulti))
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
