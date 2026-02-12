package main

import (
	"context"
	"fmt"
	"github.com/atlas/slowpoke/internal/hotel"
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

func storeProfile(ctx context.Context, req *hotel.StoreProfileRequest) *hotel.StoreProfileResponse {
	hotelId := hotel.StoreProfile(ctx, req.Profile)
	resp := hotel.StoreProfileResponse{HotelId: hotelId}
	return &resp
}

func getProfiles(ctx context.Context, req *hotel.GetProfilesRequest) *hotel.GetProfilesResponse {
	hotels := hotel.GetProfiles(ctx, req.HotelIds)
	resp := hotel.GetProfilesResponse{Profiles: hotels}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	hotel.InitProfiles()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_profile", wrappers.Wrapper[hotel.StoreProfileRequest, hotel.StoreProfileResponse](storeProfile))
	http.HandleFunc("/ro_get_profiles", wrappers.Wrapper[hotel.GetProfilesRequest, hotel.GetProfilesResponse](getProfiles))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
