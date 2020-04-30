package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
	"golang.org/x/sys/unix"
)

func init() {
	fmt.Printf("init start, os.Args = %+v\n", os.Args)
	reexec.Register("child", child)

	if reexec.Init() {
		fmt.Println("Reexec init retured non-zero")
		os.Exit(0)
	}
}

func main() {
	fmt.Println("LiteBox Starting ...")
	fmt.Printf("main start, os.Args = %+v\n", os.Args)

	mainConf := HandleFlags()

	if mainConf.Exec == "" {
		fmt.Println("No executable provided. Put executable path as --exec=<path>")
		os.Exit(1)
	}

	cmd := reexec.Command(append([]string{"child"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
		Pdeathsig:  syscall.SIGTERM,
	}

	check(cmd.Run())

}

func child() {
	fmt.Printf("child start, os.Args = %+v\t in PID: %d\n", os.Args, os.Getpid())
	conf := HandleFlags()

	cmd := exec.Command(conf.Exec)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// To drop privileges to nobody, syscall.Setuid is broken, so this workaround
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(65534), Gid: uint32(65534)},
	}
	check(syscall.Sethostname([]byte("litebox")))

	getResourceLimits()
	setResourceLimits(conf.CPU)
	getResourceLimits()
	check(cmd.Run())

}

// Not using syscall.Setrlimit it's buggy ,using unix.Setrlimit
func setResourceLimits(cpu int) {
	fmt.Println("[i] Changing resource limits")

	unix.Setrlimit(unix.RLIMIT_CPU, &unix.Rlimit{Cur: uint64(cpu), Max: uint64(cpu)})
	unix.Setrlimit(unix.RLIMIT_NOFILE, &unix.Rlimit{Cur: 50, Max: 500})
}

func getResourceLimits() {
	var getc unix.Rlimit
	var getf unix.Rlimit

	unix.Getrlimit(unix.RLIMIT_CPU, &getc)
	unix.Getrlimit(unix.RLIMIT_NOFILE, &getf)

	fmt.Println("CPU Limit: ", getc)
	fmt.Println("FILE Limit: ", getf)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
