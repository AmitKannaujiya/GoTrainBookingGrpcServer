package models

type SEAT_STATUS int

//SEAT status enum
const (
	BOOKED SEAT_STATUS = iota + 1
	AVAILABLE
)

type TRAIN_SECTION int

const (
	TRAIN_SECTION_A TRAIN_SECTION = iota + 1
	TRAIN_SECTION_B
	TRAIN_SECTION_UNDEFINED
)

type Seat struct {
	Id         int
	SeatNo     string
	Price      float64
	SeatStatus SEAT_STATUS
	Section    TRAIN_SECTION
}

func (s *Seat) UpdateSeatStatus(status SEAT_STATUS) {
	s.SeatStatus = status
}
