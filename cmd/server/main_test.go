package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/om-baji/write-go/internal"
	"github.com/om-baji/write-go/internal/utils"
	pb "github.com/om-baji/write-go/proto"
)

func setupServer(t *testing.T, dir string) *AppendServer {
	t.Helper()

	path := filepath.Join(dir, "current_0000.log")
	utils.HandlExp(utils.EnsureFile(path))

	queue := &utils.MessageQueue{}
	queue.Init(64)

	dlq := &utils.DLQueue{}
	dlq.Init(32)

	buffer := utils.NewBuffer(4096)

	seg := internal.Segment{
		Id:   1,
		Size: 0,
		Path: path,
	}

	return &AppendServer{
		queue:     queue,
		dlq:       dlq,
		buffer:    buffer,
		seg:       seg,
		threshold: 1,
	}
}

func TestAppendEnqueuesMessage(t *testing.T) {
	seqNo = 1
	dir := t.TempDir()
	srv := setupServer(t, dir)

	req := &pb.AppendRequest{
		Body:      "test-log-entry",
		Timestamp: 1718123456,
	}

	resp, err := srv.Append(context.Background(), req)
	if err != nil {
		t.Fatalf("Append returned error: %v", err)
	}

	if resp.Success != 1 {
		t.Errorf("expected success=1, got %d", resp.Success)
	}
	if resp.Error != 0 {
		t.Errorf("expected error=0, got %d", resp.Error)
	}
	if resp.Message != "append successful" {
		t.Errorf("unexpected message: %s", resp.Message)
	}

	msg := srv.queue.Dequeue()
	if msg == "" {
		t.Fatal("expected queued message, got empty string")
	}
}

func TestAppendMultipleMessages(t *testing.T) {
	seqNo = 1
	dir := t.TempDir()
	srv := setupServer(t, dir)

	for i := 0; i < 5; i++ {
		req := &pb.AppendRequest{
			Body:      fmt.Sprintf("msg-%d", i),
			Timestamp: int32(time.Now().Unix()),
		}

		resp, err := srv.Append(context.Background(), req)
		if err != nil {
			t.Fatalf("Append %d returned error: %v", i, err)
		}
		if resp.Success != 1 {
			t.Errorf("Append %d: expected success=1", i)
		}
	}

	for i := 0; i < 5; i++ {
		msg := srv.queue.Dequeue()
		if msg == "" {
			t.Fatalf("message %d should not be empty", i)
		}
	}
}

func TestProcessQueueFlushesToSegment(t *testing.T) {
	seqNo = 1
	dir := t.TempDir()
	srv := setupServer(t, dir)
	segPath := srv.seg.Path

	srv.queue.Enqueue("wal-message-one")
	srv.queue.Enqueue("wal-message-two")

	go srv.processQueue()
	time.Sleep(200 * time.Millisecond)

	content, err := os.ReadFile(segPath)
	if err != nil {
		t.Fatalf("read segment file: %v", err)
	}
	if len(content) == 0 {
		t.Error("segment file should not be empty after queue processing")
	}

	srv.queue.Enqueue("wal-message-three")
	time.Sleep(100 * time.Millisecond)

	contentAfter, _ := os.ReadFile(segPath)
	if len(contentAfter) <= len(content) {
		t.Error("segment should have grown after additional flush")
	}
}

func TestProcessDLQFlushesToSegment(t *testing.T) {
	seqNo = 1
	dir := t.TempDir()
	srv := setupServer(t, dir)
	segPath := srv.seg.Path

	srv.dlq.Enqueue("dlq-recovery-message")

	go srv.processDLQ()
	time.Sleep(200 * time.Millisecond)

	content, err := os.ReadFile(segPath)
	if err != nil {
		t.Fatalf("read segment file: %v", err)
	}
	if len(content) == 0 {
		t.Error("segment file should not be empty after dlq processing")
	}
}

func TestAppendResponseFields(t *testing.T) {
	seqNo = 1
	dir := t.TempDir()
	srv := setupServer(t, dir)

	req := &pb.AppendRequest{
		Body:      "response-test",
		Timestamp: 42,
	}

	resp, err := srv.Append(context.Background(), req)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if resp.GetSuccess() != 1 {
		t.Errorf("GetSuccess: expected 1, got %d", resp.GetSuccess())
	}
	if resp.GetMessage() != "append successful" {
		t.Errorf("GetMessage: unexpected %s", resp.GetMessage())
	}
	if resp.GetError() != 0 {
		t.Errorf("GetError: expected 0, got %d", resp.GetError())
	}
}
