package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/movie"
	"github.com/atlas/slowpoke/pkg/wrappers"
	"net/http"
	"runtime"
	"os"
	"github.com/goccy/go-json"
	"strconv"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func registerMovieId(ctx context.Context, req *movie.RegisterMovieIdRequest) *movie.RegisterMovieIdResponse {
	movie.RegisterMovieId(ctx, req.Title, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.RegisterMovieIdResponse{Ok: "OK"}
	return &resp
}

func getMovieId(ctx context.Context, req *movie.GetMovieIdRequest) *movie.GetMovieIdResponse {
	movieId := movie.GetMovieId(ctx, req.Title)
	resp := movie.GetMovieIdResponse{MovieId: movieId}
	return &resp
}
func populate() {
	ctx := context.Background()
    movie_json, err := os.ReadFile("/app/internal/movie/data/movies_1_500.json")
	if err != nil {
		panic(err)
	}
	var data []map[string]interface{}
	err = json.Unmarshal(movie_json, &data)
	if err != nil {
		panic(err)
	}
	movie_id := 0
	for _, movie_ := range data {
		movidId := strconv.Itoa(movie_id)
		title := movie_["title"].(string)
		movie.RegisterMovieId(ctx, title, movidId)
		movie_id++
	}
	fmt.Println("Populated %d movie ids", movie_id)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/register_movie_id", wrappers.Wrapper[movie.RegisterMovieIdRequest, movie.RegisterMovieIdResponse](registerMovieId))
	http.HandleFunc("/ro_get_movie_id", wrappers.Wrapper[movie.GetMovieIdRequest, movie.GetMovieIdResponse](getMovieId))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
