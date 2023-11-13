package main

import (
	"log"
	"os"

	"github.com/checkpoint-restore/go-criu/v7"
	"github.com/checkpoint-restore/go-criu/v7/rpc"
	"google.golang.org/protobuf/proto"
)

func main() {

	c := criu.MakeCriu()

	img, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("can't open image dir:", err)
	}
	defer img.Close()

	opts := &rpc.CriuOpts{
		ImagesDirFd: proto.Int32(int32(img.Fd())),
		LogLevel:    proto.Int32(4),
		ShellJob:    proto.Bool(true),
		LogFile:     proto.String("restore.log"),
	}

	c.Restore(opts, nil)
}
