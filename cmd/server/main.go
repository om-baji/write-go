package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/om-baji/write-go/internal"
	"github.com/om-baji/write-go/internal/utils"
	pb "github.com/om-baji/write-go/proto"
	"google.golang.org/grpc"
)

const WAL_MAGIC int32 = 0x57414C31

type AppendServer struct {
	pb.UnimplementedAppendServiceServer
	queue *utils.MessageQueue
	dlq   *utils.DLQueue
	seq   int64
	seg   internal.Segment
}

var seqNo int64 = 1

func (s *AppendServer) Append(
	ctx context.Context,
	req *pb.AppendRequest,
) (*pb.AppendResponse, error) {

	log.Printf(
		"received body=%s timestamp=%d",
		req.Body,
		req.Timestamp,
	)

	cr := utils.GenerateCRC([]byte(req.Body))

	entry := &internal.LedgerEntry{
		Crc:       cr,
		Body:      req.Body,
		Magic:     WAL_MAGIC,
		Timestamp: time.Now().String(),
		Seq:       seqNo,
	}

	seqNo++

	message := fmt.Sprintf("%#v", entry)
	s.queue.Enqueue(message)

	return &pb.AppendResponse{
		Success: 1,
		Message: "append successful",
		Error:   0,
	}, nil
}

func (s *AppendServer) processQueue() {
	for msg := range s.queue.GetChannel() {
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.dlq.Enqueue(msg)
				}
			}()
			s.seg = utils.AppendFlush(s.seg, msg)
		}()
	}
}

func (s *AppendServer) processDLQ() {
	for msg := range s.dlq.GetChannel() {
		func() {
			defer func() {
				if r := recover(); r != nil {
					utils.Error("dlq retry failed: %s", msg)
				}
			}()
			s.seg = utils.AppendFlush(s.seg, msg)
			preview := msg
			if len(preview) > 64 {
				preview = preview[:64]
			}
			utils.Info("dlq flush recovered: %s", preview)
		}()
	}
}

func main() {
	utils.HandlExp(
		utils.EnsureFile("./data/current_0000.log"),
	)

	queue := &utils.MessageQueue{}
	queue.Init(1024)

	dlq := &utils.DLQueue{}
	dlq.Init(256)

	seg := internal.Segment{
		Id:   1,
		Size: 0,
		Path: "./data/current_0000.log",
	}

	server := &AppendServer{
		queue: queue,
		dlq:   dlq,
		seg:   seg,
	}

	go server.processQueue()
	go server.processDLQ()

	lis, err := net.Listen("tcp", ":50051")
	utils.HandlExp(err)

	grpcServer := grpc.NewServer()

	pb.RegisterAppendServiceServer(grpcServer, server)

	utils.Info("Starting Server...")
	utils.Info("Listening at 50051")

	if err = grpcServer.Serve(lis); err != nil {
		utils.Error("Something went wrong %v", err.Error())
	}
}
