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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	check(cmd.Run())

}

func child() {
	fmt.Printf("Running %s in PID: %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// To drop privileges to nobody, syscall.Setuid is broken, so this workaround
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(65534), Gid: uint32(65534)},
	}
	check(syscall.Sethostname([]byte("litebox")))

	getResourceLimits()
	setResourceLimits()
	getResourceLimits()
	check(cmd.Run())

}

func setResourceLimits() {
	fmt.Println("Changing resource limits")
	//unix.Setrlimit(unix.RLIMIT_CPU, &unix.Rlimit{Cur: uint64(5), Max: uint64(5)})
	syscall.Setrlimit(syscall.RLIMIT_CPU, &syscall.Rlimit{Cur: 5, Max: 500})
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 50, Max: 500})
}

func getResourceLimits() {
	var getc syscall.Rlimit
	var getf syscall.Rlimit

	syscall.Getrlimit(syscall.RLIMIT_CORE, &getc)
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &getf)
	fmt.Println("CPU Limit: ", getc)
	fmt.Println("FILE Limit: ", getf)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
