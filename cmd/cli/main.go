package main

import (
	"fmt"
	"os"
	"time"

	"github.com/om-baji/write-go/internal"
	"github.com/om-baji/write-go/internal/utils"
)

func init() {
	utils.HandlExp(
		utils.EnsureFile(
			"./data/current.log",
		),
	)
}

const WAL_MAGIC int32 = 0x57414C31

var Seq_No = 1

var CurrentSegment = internal.Segment{
	Id:   1,
	Size: 0,
	Path: "./data/current_0000.log",
}

func help() {
	fmt.Println("A go based WAL engine with gRPC support.")
	fmt.Println("Example :")
	fmt.Println("\tcli add <log>\t Add a log")
	fmt.Println("\tcli verify\t verifies the current logs")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	command := os.Args[1]

	switch command {
	case "add":
		payload := os.Args[2]
		cr := utils.GenerateCRC([]byte(payload))

		entry := &internal.LedgerEntry{
			Crc:       cr,
			Body:      payload,
			Magic:     WAL_MAGIC,
			Timestamp: time.Now().String(),
			Seq:       int64(Seq_No),
		}

		message := fmt.Sprintf("%#v", entry)

		CurrentSegment = utils.AppendFlush(CurrentSegment, message)
	case "verify":
		println("This is add command!")
	case "info":
		println("id : ", CurrentSegment.Id)
		println("path : ", CurrentSegment.Path)
		println("size : ", CurrentSegment.Size)
	default:
		println("Invalid Command!")
	}
}
