package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	// Process conf.Exec
	c := strings.Split(conf.Exec, " ")
	exe := c[0]
	args := c[1:]

	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// To drop privileges to nobody, syscall.Setuid is broken, so this workaround
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(65534), Gid: uint32(65534)},
	}
	check(syscall.Sethostname([]byte("litebox")))

	// getResourceLimits()
	setResourceLimits(conf.CPU, conf.Memory, conf.Nproc)
	getResourceLimits()

	check(cmd.Run())

	// Get memory usage after process exits
	usage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	fmt.Println("Memory Used: ", usage.Maxrss)
	// Calutale Seconds in floating point
	cpuTime := float64(usage.Stime.Usec+usage.Utime.Usec) / float64(100000) // No idea what is Usec
	fmt.Println("CPU time(usr+sys): ", cpuTime)
}

// Not using syscall.Setrlimit it's buggy ,using unix.Setrlimit
func setResourceLimits(cpu int, mem int, nproc int) {
	fmt.Println("[i] Changing resource limits")

	unix.Setrlimit(unix.RLIMIT_CPU, &unix.Rlimit{Cur: uint64(cpu), Max: uint64(cpu)})
	unix.Setrlimit(unix.RLIMIT_AS, &unix.Rlimit{Cur: uint64(mem), Max: uint64(mem)})
	unix.Setrlimit(unix.RLIMIT_NPROC, &unix.Rlimit{Cur: uint64(nproc), Max: uint64(nproc)})
}

func getResourceLimits() {
	var getc unix.Rlimit
	var getn unix.Rlimit
	var getm unix.Rlimit

	unix.Getrlimit(unix.RLIMIT_CPU, &getc)
	unix.Getrlimit(unix.RLIMIT_AS, &getm)
	unix.Getrlimit(unix.RLIMIT_NPROC, &getn)

	fmt.Println("CPU Limit: ", getc)
	fmt.Println("Memory Limit: ", getm)
	fmt.Println("Process Limit: ", getn)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
