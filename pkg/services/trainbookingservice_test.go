package services

import (
	"testing"
	conf "train-book/cmd/config"
	pb "train-book/proto/api.v1"

	"github.com/stretchr/testify/assert"
)

func TestCreateBooking(t *testing.T) {
	service := GetTrainBookingServiceInstance()
	config, _ := conf.GetConfig()
	service.CreateAndLoadSeats(config)

    req := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
        From: "London", To: "France", Price: 20.0, Section: "A",
    }

    res, err := service.CreateBooking(req)
    assert.Nil(t, err)
    assert.NotEmpty(t, res.BookingId)
}

func TestGetBookingStatus(t *testing.T) {
	service := GetTrainBookingServiceInstance()
	config, _ := conf.GetConfig()
	service.CreateAndLoadSeats(config)
	req := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
        From: "London", To: "France", Price: 20.0, Section: "A",
    }

    res1, _ := service.CreateBooking(req)

	req2 := &pb.ReceiptRequest{ReceiptId: res1.BookingId}
    res2, err := service.GetBookingStatus(req2.ReceiptId)
    assert.Nil(t, err)
    assert.Equal(t, "London", res2.Source)
	assert.Equal(t, "France", res2.Destination)
	//assert.Equal(t, "john.doe@example.com", booking.User.Email)
}

func TestRemoveUserFromBooking(t *testing.T) {
	service := GetTrainBookingServiceInstance()
	config, _ := conf.GetConfig()
	service.CreateAndLoadSeats(config)
	req := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "Jane", LastName: "Doe", Email: "Jane.doe@example.com"},
        From: "London", To: "France", Price: 10.0, Section: "B",
    }

    _,  err := service.CreateBooking(req)

	assert.Nil(t, err)

    req2 := &pb.RemoveRequest{Email: req.User.Email}
    err1 := service.RemoveUserFromBooking(req2)
    assert.Nil(t, err1)
	assert.Equal(t, "Jane.doe@example.com", req.User.Email)
}

func TestModifyUserBookedSeat(t *testing.T) {
	service := GetTrainBookingServiceInstance()
	config, _ := conf.GetConfig()
	service.CreateAndLoadSeats(config)
	req1 := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "Alice", LastName: "Smith", Email: "alice.smith@example.com"},
        From: "London", To: "France", Price: 10.0, Section: "B",
    }

    res1, err1 := service.CreateBooking(req1)
	assert.Nil(t, err1)
    req2 := &pb.ModifySeatRequest{Email: req1.User.Email, NewSeat: "A5"}
    err2 := service.ModifyUserBookedSeat(req2)
    assert.Nil(t, err2)
	updatedBooking, err3 := service.GetBookingStatus(res1.BookingId)
	assert.Nil(t, err3)
    assert.Equal(t, "A5", updatedBooking.Seat.SeatNo)
}

func TestGetBookedUserSeatBySection(t *testing.T) {
	service := GetTrainBookingServiceInstance()
	config, _ := conf.GetConfig()
	service.CreateAndLoadSeats(config)
	req1 := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
        From: "London", To: "France", Price: 10.0, Section: "A",
    }

    res1, err1 := service.CreateBooking(req1)
	assert.Nil(t, err1)

	req2 := &pb.PurchaseRequest{
        User: &pb.User{FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com"},
        From: "Paris", To: "Berlin", Price: 10.0, Section: "B",
    }

    res2, err2 := service.CreateBooking(req2)
	assert.Nil(t, err2)

    // Testing for Section A
    reqA := &pb.SectionRequest{Section: "A"}
    resA := service.GetBookedUserSeatBySection(reqA.Section)
    assert.Equal(t, 1, len(resA))
    assert.Equal(t, res1.Seat.SeatNo, resA[0].Seat.SeatNo)
    assert.Equal(t, "John", resA[0].User.FirstName)
	//assert.Equal(t, "receipt-4", booking1.BookingId)

    // Testing for Section B
    reqB := &pb.SectionRequest{Section: "B"}
    resB := service.GetBookedUserSeatBySection(reqB.Section)
    assert.Equal(t, 1, len(resB))
    assert.Equal(t, res2.Seat.SeatNo, resB[0].Seat.SeatNo)
    assert.Equal(t, "Jane", resB[0].User.FirstName)
	//assert.Equal(t, "receipt-5", booking2.BookingId)
}
