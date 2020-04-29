package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println("LiteBox Starting ...")

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("An Error occured")
	}
}

func run() {
	fmt.Printf("Running %s in PID: %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	check(cmd.Run())

}

func child() {
	fmt.Printf("Running %s in PID: %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Drop privilages to nobody, syscall.Setuid is broken, so this workaround
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(65534), Gid: uint32(65534)},
	}

	check(syscall.Sethostname([]byte("litebox")))
	check(cmd.Run())

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
