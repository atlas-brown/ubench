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

func uploadMovieReview(ctx context.Context, req *movie.UploadMovieReviewRequest) *movie.UploadMovieReviewResponse {
	reviewId := movie.UploadMovieReview(ctx, req.MovieId, req.ReviewId, req.Timestamp)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.UploadMovieReviewResponse{ReviewId: reviewId}
	return &resp
}

func readMovieReviews(ctx context.Context, req *movie.ReadMovieReviewsRequest) *movie.ReadMovieReviewsResponse {
	reviews := movie.ReadMovieReviews(ctx, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadMovieReviewsResponse{Reviews: reviews}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/upload_movie_review", wrappers.Wrapper[movie.UploadMovieReviewRequest, movie.UploadMovieReviewResponse](uploadMovieReview))
	http.HandleFunc("/ro_read_movie_reviews", wrappers.Wrapper[movie.ReadMovieReviewsRequest, movie.ReadMovieReviewsResponse](readMovieReviews))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
