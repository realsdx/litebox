package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

// ShowUsage shows the amount of time the CPU(sys+usr) was in use for the program
func ShowUsage(cmd *exec.Cmd) {
	// Get memory usage after process exits
	usage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	fmt.Printf("Memory Used(RSS): %d kB\n", usage.Maxrss)
	// Calutale Seconds in floating point
	cpuTimeMicro := (usage.Stime.Sec*1000000 + usage.Stime.Usec) + (usage.Utime.Sec*1000000 + usage.Utime.Usec)

	fmt.Printf("CPU time(usr+sys): %d MicroSec\n", cpuTimeMicro)
	cpuTime := float64(cpuTimeMicro) / float64(1000000)
	fmt.Printf("CPU time(usr+sys): %f Sec\n", cpuTime)
}
