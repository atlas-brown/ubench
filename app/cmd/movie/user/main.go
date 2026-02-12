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

func registerUser(ctx context.Context, req *movie.RegisterUserRequest) *movie.RegisterUserResponse {
	ok := movie.RegisterUser(ctx, req.Username, req.Password)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.RegisterUserResponse{Ok: ok}
	return &resp
}

func login(ctx context.Context, req *movie.LoginRequest) *movie.LoginResponse {
	token := movie.Login(ctx, req.Username, req.Password)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.LoginResponse{Token: token}
	return &resp
}

func getUserId(ctx context.Context, req *movie.GetUserIdRequest) *movie.GetUserIdResponse {
	userId := movie.GetUserId(ctx, req.Username)
	resp := movie.GetUserIdResponse{UserId: userId}
	return &resp
}

func populate() {
	for i := 0; i < 100; i++ {
		movie.RegisterUser(context.Background(), fmt.Sprintf("username%d", i), fmt.Sprintf("password%d", i))
	}
	fmt.Println("Populated %d users", 100)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/register_user", wrappers.Wrapper[movie.RegisterUserRequest, movie.RegisterUserResponse](registerUser))
	http.HandleFunc("/login", wrappers.Wrapper[movie.LoginRequest, movie.LoginResponse](login))
	http.HandleFunc("/ro_get_user_id", wrappers.Wrapper[movie.GetUserIdRequest, movie.GetUserIdResponse](getUserId))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
