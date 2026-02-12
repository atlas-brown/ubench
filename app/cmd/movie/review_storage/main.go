package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/movie"
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

func storeReview(ctx context.Context, req *movie.StoreReviewRequest) *movie.StoreReviewResponse {
	reviewId := movie.StoreReview(ctx, req.Review)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.StoreReviewResponse{ReviewId: reviewId}
	return &resp
}

func readReviews(ctx context.Context, req *movie.ReadReviewsRequest) *movie.ReadReviewsResponse {
	reviews := movie.ReadReviews(ctx, req.ReviewIds)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadReviewsResponse{Reviews: reviews}
	//fmt.Printf("[ReviewStorage] Response: %v\n", resp)
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_review", wrappers.Wrapper[movie.StoreReviewRequest, movie.StoreReviewResponse](storeReview))
	http.HandleFunc("/ro_read_reviews", wrappers.Wrapper[movie.ReadReviewsRequest, movie.ReadReviewsResponse](readReviews))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
