package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "train-book/proto/api.v1"
	u "train-book/utill"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedComServiceServer
	pb.UnimplementedTrainBookingServer
	bookingMap map[string]string
}

func (s *Server) ProcessMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("Received message: %v", req)

	return &pb.Response{
		Name: req.Name,
		ClientMessage: req.Message,
		ServerMessage: "Hello from grpc server",
		ParentId: req.Id,
		ChildId: req.Id * 3,
	}, nil
}

func (s *Server) BookTicket(ctx context.Context, req *pb.TicketRequest) (*pb.TicketResponse, error) {
	log.Printf("Booking Received message: %v", req)
	confirmation := fmt.Sprintf("conf-%d", len(s.bookingMap) + 1)
	s.bookingMap[confirmation] = "Booked"
	return &pb.TicketResponse{
		Confirmation: confirmation,
		Message: "Ticket successfully booked :\n Mr. " + req.FirstName + " " +  req.LastName + "\nEmail : " + req.Email +
		"\n From : " + req.Departure + " to : "+ req.Destination +  " for $ : " + fmt.Sprintf("%.2f", req.Price),
	}, nil
}

func (s *Server) GetBookingStatus(ctx context.Context, req *pb.BookingStatusRequest) (*pb.BookingStatusResponse, error) {
	log.Printf("Booking Status Received message: %v", req)

	status , ok := s.bookingMap[req.Confirmation]
	if !ok {
		return &pb.BookingStatusResponse{Status: "Not found"}, nil
	} 
	return &pb.BookingStatusResponse{Status: status}, nil
}

func main() {
	fmt.Println("Hello From Server")
	listen, err := net.Listen("tcp", u.PORT)

	if err != nil {
		fmt.Println("Error happend : %v", err)
	}
	s := grpc.NewServer()
	sr := &Server{
		bookingMap: make(map[string]string),
	}
	pb.RegisterComServiceServer(s, sr)
	pb.RegisterTrainBookingServer(s, sr)
	log.Printf("server listening on %v", u.PORT)
	if er := s.Serve(listen); er != nil {
		log.Fatalf("failed to serve : %v", er)
	}
}