package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	criu "github.com/checkpoint-restore/go-criu/v7"
	"github.com/checkpoint-restore/go-criu/v7/rpc"
	"google.golang.org/protobuf/proto"
)

type NoNotify struct {
	criu.NoNotify
}

// steps:
// setsid go run /home/hannahmanuela/gotest/sleep/sleep.go < /dev/null &> sleep_out.log &
// ps -aux | less | grep sleep.go
// --> put that pid into sudo go run test.go <pid>

func CheckPoint(imgDir string, pid int, c *criu.Criu) error {

	img, err := os.Open(imgDir)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("can't open image dir")
	}
	defer img.Close()

	opts := &rpc.CriuOpts{
		Pid:         proto.Int32(int32(pid)),
		ImagesDirFd: proto.Int32(int32(img.Fd())),
		LogLevel:    proto.Int32(4),
		ShellJob:    proto.Bool(true),
		LogFile:     proto.String("dump.log"),
	}

	return c.Dump(opts, NoNotify{})

}

func Restore(imgDir string, c *criu.Criu) error {

	img, err := os.Open(imgDir)
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

	return c.Restore(opts, nil)
}

func main() {

	imgDir := "chkptimg"
	c := criu.MakeCriu()

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("invalid pid " + os.Args[1])
	}

	err = CheckPoint(imgDir, pid, c)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("dump fail")
	} else {
		fmt.Println("dump successful")
	}
	// wait
	time.Sleep(2 * time.Second)

	fmt.Println("trying to restore")

	// restore
	err = Restore(imgDir, c)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("restore fail")
	} else {
		fmt.Println("restore successful")
	}

}
