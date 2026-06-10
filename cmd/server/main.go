package main

import (
	"context"
	"log"
	"net"

	"github.com/om-baji/write-go/internal/utils"
	pb "github.com/om-baji/write-go/proto"
	"google.golang.org/grpc"
)

type AppendServer struct {
	pb.UnimplementedAppendServiceServer
}

var queue *utils.MessageQueue

func (s *AppendServer) Append(
	ctx context.Context,
	req *pb.AppendRequest,
) (*pb.AppendResponse, error) {

	log.Printf(
		"received body=%s timestamp=%d",
		req.Body,
		req.Timestamp,
	)
	queue.Enqueue(req.Body)

	return &pb.AppendResponse{
		Success: 1,
		Message: "append successful",
		Error:   0,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	utils.HandlExp(err)

	server := grpc.NewServer()

	pb.RegisterAppendServiceServer(server, &AppendServer{})

	utils.Info("Starting Server...")
	utils.Info("Listening at 50001")

	if err = server.Serve(lis); err != nil {
		utils.Error("Something went wrong %v", err.Error())
	}
}
