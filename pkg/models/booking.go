package models

type BOOKING_STATUS int

const (
	BOOKING_STATUS_CONFIRMED BOOKING_STATUS = iota + 1
	BOOKING_STATUS_CANCELLED
	BOOKING_STATUS_FAILED
	BOOKING_STATUS_PENDING
)

type Booking struct {
	BookingId     string
	User          User
	Source        string
	Destination   string
	BookingDate   string
	BookingAt     string
	BookingStatus BOOKING_STATUS
	Seat          Seat
	Price         float64
}

type UserSeat struct {
	User User
	Seat Seat
}

func (b *Booking) UpdateBookingStatus(status BOOKING_STATUS) {
	b.BookingStatus = status
}

func (b *Booking) ModifySeat(seat Seat) {
	b.Seat = seat
}
