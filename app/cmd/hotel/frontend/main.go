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

func searchHotels(ctx context.Context, req *hotel.SearchHotelsRequest) *hotel.SearchHotelsResponse {
	hotels := hotel.SearchHotels(ctx, req.InDate, req.OutDate, req.Location)
	resp := hotel.SearchHotelsResponse{Profiles: hotels}
	return &resp
}

func storeHotel(ctx context.Context, req *hotel.StoreHotelRequest) *hotel.StoreHotelResponse {
	hotelId := hotel.StoreHotel(ctx, req.HotelId, req.Name, req.Phone, req.Location, req.Rate, req.Capacity, req.Info)
	resp := hotel.StoreHotelResponse{HotelId: hotelId}
	return &resp
}

func reservation(ctx context.Context, req *hotel.FrontendReservationRequest) *hotel.FrontendReservationResponse {
	success := hotel.FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, req.Rooms, req.Username, req.Password)
	resp := hotel.FrontendReservationResponse{Success: success}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_search_hotels", wrappers.Wrapper[hotel.SearchHotelsRequest, hotel.SearchHotelsResponse](searchHotels))
	http.HandleFunc("/store_hotel", wrappers.Wrapper[hotel.StoreHotelRequest, hotel.StoreHotelResponse](storeHotel))
	http.HandleFunc("/reservation", wrappers.Wrapper[hotel.FrontendReservationRequest, hotel.FrontendReservationResponse](reservation))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
