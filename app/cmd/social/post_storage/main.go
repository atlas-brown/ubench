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

func storePost(ctx context.Context, req *social.StorePostRequest) *social.StorePostResponse {
	postId := social.StorePost(ctx, req.CreatorId, req.Text)
	//fmt.Println("Post stored: " + postId)
	resp := social.StorePostResponse{PostId: postId}
	return &resp
}

func storePostMulti(ctx context.Context, req *social.StorePostMultiRequest) *social.StorePostMultiResponse {
	postIds := social.StorePostMulti(ctx, req.CreatorId, req.Text, req.Number)
	//fmt.Println("Post stored: " + postId)
	resp := social.StorePostMultiResponse{PostIds: postIds}
	return &resp
}

func readPost(ctx context.Context, req *social.ReadPostRequest) *social.ReadPostResponse {
	post := social.ReadPost(ctx, req.PostId)
	//fmt.Printf("Post read: %+v\n", post)
	resp := social.ReadPostResponse{Post: post}
	return &resp
}

func readPosts(ctx context.Context, req *social.ReadPostsRequest) *social.ReadPostsResponse {
	posts := social.ReadPosts(ctx, req.PostIds)
	//fmt.Printf("Posts read: %+v\n", posts)
	resp := social.ReadPostsResponse{Posts: posts}
	return &resp
}

func main() {
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_post", wrappers.Wrapper[social.StorePostRequest, social.StorePostResponse](storePost))
	http.HandleFunc("/store_post_multi", wrappers.Wrapper[social.StorePostMultiRequest, social.StorePostMultiResponse](storePostMulti))
	http.HandleFunc("/ro_read_post", wrappers.Wrapper[social.ReadPostRequest, social.ReadPostResponse](readPost))
	http.HandleFunc("/ro_read_posts", wrappers.Wrapper[social.ReadPostsRequest, social.ReadPostsResponse](readPosts))
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
