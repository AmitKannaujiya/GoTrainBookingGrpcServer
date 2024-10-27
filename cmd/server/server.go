package server

import (
	"context"
	"log"
	"net"
	pb "train-book/proto/api.v1"

	"google.golang.org/grpc"
	"sync"
	cnf "train-book/cmd/config"
	service "train-book/pkg/services"
)

type Server struct {
	// start from here
	TrainBookingService service.TrainBookingService
	pb.UnimplementedTrainTicketingServer
}

var serverInstance *Server
var once sync.Once

func GetServerInstance() *Server {
	if serverInstance == nil {
		once.Do(func() {
			serverInstance = &Server{
				TrainBookingService: *service.GetTrainBookingServiceInstance(),
			}
		})
	}
	return serverInstance
}

func (s *Server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	log.Printf("Received message: %v", req)
	booking, err := s.TrainBookingService.CreateBooking(req)
	if err != nil {
		log.Printf("Error happend : %v", err)
		return &pb.PurchaseResponse{Success: false, Detail: err.Error()}, err
	}
	return &pb.PurchaseResponse{
		Success:   true,
		User:      req.User,
		ReceiptId: booking.BookingId,
		Seat:      booking.Seat.SeatNo,
		Detail:    "Mr. " + booking.User.GetFullName() + " seat is booked Successfull",
	}, nil
}

func (s *Server) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
	log.Printf("Received message: %v", req)
	booking, err := s.TrainBookingService.GetBookingStatus(req.ReceiptId)
	if err != nil {
		log.Printf("Error happend : %v", err)
		return &pb.ReceiptResponse{
			Success: false,
		}, err
	}
	return &pb.ReceiptResponse{
		Success: true,
		From:    booking.Source,
		To:      booking.Destination,
		User: &pb.User{
			FirstName: booking.User.FirstName,
			LastName:  booking.User.LastName,
			Email:     booking.User.Email,
		},
		Seat:  booking.Seat.SeatNo,
		Price: booking.Price,
	}, nil
}

func (s *Server) GetUsersBySection(ctx context.Context, req *pb.SectionRequest) (*pb.UsersResponse, error) {
	log.Printf("Received message: %v", req)
	userSeatMapping := s.TrainBookingService.GetBookedUserSeatBySection(req.Section)
	var userSeats []*pb.UserSeat
	for _, userSeat := range userSeatMapping {
		userSeats = append(userSeats,
			&pb.UserSeat{
				User: &pb.User{
					FirstName: userSeat.User.FirstName,
					LastName:  userSeat.User.LastName,
					Email:     userSeat.User.Email,
				},
				Seat: userSeat.Seat.SeatNo,
			})
	}
	return &pb.UsersResponse{
		Users: userSeats,
	}, nil
}
func (s *Server) RemoveUser(ctx context.Context, req *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	log.Printf("Received message: %v", req)
	err := s.TrainBookingService.RemoveUserFromBooking(req)
	if err != nil {
		log.Printf("Error happend : %v", err)
		return &pb.RemoveResponse{
			Success: false,
		}, err
	}
	return &pb.RemoveResponse{
		Success: true,
	}, nil
}

func (s *Server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	log.Printf("Received message: %v", req)
	err := s.TrainBookingService.ModifyUserBookedSeat(req)
	if err != nil {
		log.Printf("Error happend : %v", err)
		return &pb.ModifySeatResponse{
			Success: false,
		}, err
	}
	return &pb.ModifySeatResponse{
		Success: true,
	}, nil
}

func Execute(Config *cnf.Config) {
	address := Config.App.Host + ":" + Config.App.Port
	listen, err := net.Listen("tcp", address)

	if err != nil {
		log.Printf("Error happend : %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	server := GetServerInstance()
	pb.RegisterTrainTicketingServer(grpcServer, server)
	log.Printf("server listening on %v", address)
	if er := grpcServer.Serve(listen); er != nil {
		log.Printf("failed to serve : %v", er)
	}
}
