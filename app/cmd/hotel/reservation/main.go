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

func checkAvailability(ctx context.Context, req *hotel.CheckAvailabilityRequest) *hotel.CheckAvailabilityResponse {
	hotelIds := hotel.CheckAvailability(ctx, req.CustomerName, req.HotelIds, req.InDate, req.OutDate, req.RoomNumber)
	resp := hotel.CheckAvailabilityResponse{HotelIds: hotelIds}
	return &resp
}

func makeReservation(ctx context.Context, req *hotel.MakeReservationRequest) *hotel.MakeReservationResponse {
	success := hotel.MakeReservation(ctx, req.CustomerName, req.HotelId, req.InDate, req.OutDate, req.RoomNumber)
	resp := hotel.MakeReservationResponse{Success: success}
	return &resp
}

func addHotelAvailability(ctx context.Context, req *hotel.AddHotelAvailabilityRequest) *hotel.AddHotelAvailabilityResponse {
	hotelId := hotel.AddHotelAvailability(ctx, req.HotelId, req.Capacity)
	resp := hotel.AddHotelAvailabilityResponse{Hotelid: hotelId}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	hotel.InitHotelAvailability()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_check_availability", wrappers.Wrapper[hotel.CheckAvailabilityRequest, hotel.CheckAvailabilityResponse](checkAvailability))
	http.HandleFunc("/make_reservation", wrappers.Wrapper[hotel.MakeReservationRequest, hotel.MakeReservationResponse](makeReservation))
	http.HandleFunc("/add_hotel_availability", wrappers.Wrapper[hotel.AddHotelAvailabilityRequest, hotel.AddHotelAvailabilityResponse](addHotelAvailability))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
