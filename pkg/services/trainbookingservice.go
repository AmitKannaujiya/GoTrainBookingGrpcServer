package services

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"
	m "train-book/pkg/models"
	"train-book/utill"

	c "train-book/cmd/config"
	pb "train-book/proto/api.v1"
	u "train-book/utill"

	"github.com/google/uuid"
)

var once = &sync.Once{}
var lock = &sync.RWMutex{}

type TrainBookingService struct {
	Users       map[string]m.User
	Seats       []*m.Seat
	BookedSeats map[string]*m.Booking
	BookedById  map[string]*m.Booking
}

var SingletonInstance *TrainBookingService

func GetTrainBookingServiceInstance() *TrainBookingService {
	if SingletonInstance == nil {
		once.Do(func() {
			SingletonInstance = &TrainBookingService{
				Users:       make(map[string]m.User),
				Seats:       make([]*m.Seat, 0),
				BookedSeats: make(map[string]*m.Booking),
				BookedById:  make(map[string]*m.Booking),
			}
		})
	}
	return SingletonInstance
}

// initilize seat in train
func (tb *TrainBookingService) CreateAndLoadSeats(Config *c.Config) {
	lock.Lock()
	defer lock.Unlock()
	seatPerSection := Config.Train.SeatPerSection
	for i := 0; i < seatPerSection; i++ {
		seatA := &m.Seat{
			Id:         i + 1,
			SeatNo:     utill.SECTION_A + strconv.Itoa(i+1),
			Section:    m.TRAIN_SECTION_A,
			SeatStatus: m.AVAILABLE,
		}
		tb.Seats = append(tb.Seats, seatA)
		seatB := &m.Seat{
			Id:         i + 1,
			SeatNo:     utill.SECTION_B + strconv.Itoa(i+1),
			Section:    m.TRAIN_SECTION_B,
			SeatStatus: m.AVAILABLE,
		}
		tb.Seats = append(tb.Seats, seatB)
	}
}

func (tb *TrainBookingService) CreateBooking(req *pb.PurchaseRequest) (*m.Booking, error) {
	lock.Lock()
	defer lock.Unlock()
	firstName, lastName, email, section, to, from, date := req.User.FirstName, req.User.LastName, req.User.Email, u.GetSectionType(req.Section), req.To, req.From, req.When
	user, exists := tb.Users[email]
	if !exists {
		user = m.User{
			Id:        len(tb.Users) + 1,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}
		tb.Users[email] = user
	}

	seat, er := tb.GetNextAvailableSeat(user, date, from, to, section)
	if er != nil {
		log.Printf("No seat available : %v", er)
		return nil, er
	}
	seat.UpdateSeatStatus(m.BOOKED)
	today := time.Now().Format("2017-09-21")
	uuid := uuid.New()
	booking := &m.Booking{
		BookingId:     uuid.String(),
		User:          user,
		Source:        from,
		Destination:   to,
		BookingDate:   date,
		BookingStatus: m.BOOKING_STATUS_CONFIRMED,
		Seat:          *seat,
		BookingAt:     today,
	}
	tb.BookedSeats[email] = booking
	tb.BookedById[uuid.String()] = booking
	return booking, nil
}

func (tb *TrainBookingService) GetNextAvailableSeat(user m.User, date, to, from string, section m.TRAIN_SECTION) (*m.Seat, error) {
	if booking, exists := tb.BookedSeats[user.Email]; exists {
		if booking.BookingDate == date && booking.Source == to && booking.Destination == from {
			return nil, errors.New("BOOKING ALREADY PRESENT")
		}
	}
	for _, seat := range tb.Seats {
		if seat.Section == section && seat.SeatStatus == m.AVAILABLE {
			return seat, nil
		}
	}
	return nil, errors.New("NO SEATS AVAILABLE")
}

func (tb *TrainBookingService) GetBookingStatus(bookingId string) (*m.Booking, error) {
	lock.Lock()
	defer lock.Unlock()
	if booking, exists := tb.BookedById[bookingId]; exists {
		return booking, nil
	}
	return nil, errors.New("NO SUCH BOOKING EXISTS")
}

func (tb *TrainBookingService) GetBookedUserSeatBySection(section string) []m.UserSeat {
	lock.Lock()
	defer lock.Unlock()
	requestedSection := u.GetSectionType(section)
	var userSeatDetails []m.UserSeat
	for _, booking := range tb.BookedSeats {
		if booking.Seat.Section == requestedSection {
			userSeatDetails = append(userSeatDetails, m.UserSeat{
				User: booking.User,
				Seat: booking.Seat,
			})
		}
	}
	return userSeatDetails
}

func (tb *TrainBookingService) RemoveUserFromBooking(req *pb.RemoveRequest) error {
	lock.Lock()
	defer lock.Unlock()
	email := req.Email
	if booking, exists := tb.BookedSeats[email]; exists {
		booking.UpdateBookingStatus(m.BOOKING_STATUS_FAILED)
		booking.Seat.UpdateSeatStatus(m.AVAILABLE)
		for _, seat := range tb.Seats {
			if seat.Id == booking.Seat.Id {
				seat.UpdateSeatStatus(m.AVAILABLE)
			}
		}
		return nil
	} else {
		return errors.New("NO BOOKING FOUND")
	}
}

func (tb *TrainBookingService) ModifyUserBookedSeat(req *pb.ModifySeatRequest) error {
	lock.Lock()
	defer lock.Unlock()
	email, newSeatNo := req.Email, req.NewSeat
	if booking, exists := tb.BookedSeats[email]; exists {
		oldSeat := booking.Seat
		var newSeat m.Seat
		for _, seat := range tb.Seats {
			if seat.SeatNo == newSeatNo {
				newSeat = *seat
				break
			}
		}
		if newSeat.SeatStatus == m.AVAILABLE {
			oldSeat.UpdateSeatStatus(m.AVAILABLE)
			newSeat.UpdateSeatStatus(m.BOOKED)
			booking.Seat = newSeat
			return nil
		} else {
			return errors.New("SEAT ALREADY BOOKED")
		}
	} else {
		return errors.New("BOOKING NOT AVAILABLE")
	}
}
