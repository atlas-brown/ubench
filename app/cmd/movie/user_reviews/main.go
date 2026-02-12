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

func uploadUserReview(ctx context.Context, req *movie.UploadUserReviewRequest) *movie.UploadUserReviewResponse {
	reviewId := movie.UploadUserReview(ctx, req.UserId, req.ReviewId, req.Timestamp)
	//fmt.Println("User info stored for id: " + movieId)
	resp := movie.UploadUserReviewResponse{ReviewId: reviewId}
	return &resp
}

func readUserReviews(ctx context.Context, req *movie.ReadUserReviewsRequest) *movie.ReadUserReviewsResponse {
	reviews := movie.ReadUserReviews(ctx, req.UserId)
	//fmt.Printf("User info read: %v\n", movieInfo)
	resp := movie.ReadUserReviewsResponse{Reviews: reviews}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/upload_user_review", wrappers.Wrapper[movie.UploadUserReviewRequest, movie.UploadUserReviewResponse](uploadUserReview))
	http.HandleFunc("/ro_read_user_reviews", wrappers.Wrapper[movie.ReadUserReviewsRequest, movie.ReadUserReviewsResponse](readUserReviews))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
