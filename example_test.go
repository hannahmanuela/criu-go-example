package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"

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
//
// to not need to run sudo, grant binary capabilities:
// sudo setcap cap_sys_ptrace,cap_sys_admin,cap_dac_read_search,cap_net_admin=eip /usr/local/sbin/criu

func CheckPoint(imgDir string, pid int, c *criu.Criu) error {

	img, err := os.Open(imgDir)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("can't open image dir")
	}
	defer img.Close()

	opts := &rpc.CriuOpts{
		Pid:          proto.Int32(int32(pid)),
		ImagesDirFd:  proto.Int32(int32(img.Fd())),
		LogLevel:     proto.Int32(4),
		ShellJob:     proto.Bool(true),
		Unprivileged: proto.Bool(true),
		LogFile:      proto.String("dump.log"),
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
		ImagesDirFd:  proto.Int32(int32(img.Fd())),
		LogLevel:     proto.Int32(4),
		ShellJob:     proto.Bool(true),
		Unprivileged: proto.Bool(true),
		LogFile:      proto.String("restore.log"),
	}

	return c.Restore(opts, nil)
}

func TestCheckpointing(t *testing.T) {

	imgDir := "chkptimg"
	c := criu.MakeCriu()

	// log.Printf(os.Getenv("PATH"))
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/home/hannahmanuela/example-criu-go")
	// log.Printf(os.Getenv("PATH"))

	cmd := exec.Command("example")

	outfile, err := os.Create("./out.txt")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	// cmd.Stdout = os.Stdout

	// Set up new namespaces
	// cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error start %v %v", cmd, err)
	}
	log.Printf("---> RUNNING WITH PID %d\n", cmd.Process.Pid)

	err = CheckPoint(imgDir, cmd.Process.Pid, c)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("dump fail")
	} else {
		fmt.Println("dump successful")
	}
}

// func TestRestoring(t *testing.T) {
// 	fmt.Println("trying to restore")

// 	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/home/hannahmanuela/example-criu-go")

// 	Restore("chkptimg", criu.MakeCriu())

// 	// restore
// 	// cmd := exec.Command("restore-wrapper", "./chkptimg")
// 	// cmd.Stdout = os.Stdout
// 	// cmd.Stderr = os.Stderr

// 	// if err := cmd.Start(); err != nil {
// 	// log.Fatalf("Error start %v %v", cmd, err)
// 	// }
// 	// log.Printf("---> RUNNING WITH PID %d\n", cmd.Process.Pid)
// }
