package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// ShowCPUTime shows the amount of time the CPU(sys+usr) was in use for the program
func ShowCPUTime(cmd *exec.Cmd) {
	// Get memory usage after process exits
	usage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	fmt.Println("Memory Used(RSS): ", usage.Idrss, usage.Isrss)
	// Calutale Seconds in floating point
	cpuTimeMicro := (usage.Stime.Sec*1000000 + usage.Stime.Usec) + (usage.Utime.Sec*1000000 + usage.Utime.Usec)

	fmt.Printf("CPU time(usr+sys): %d Usec\n", cpuTimeMicro)
	cpuTime := float64(cpuTimeMicro) / float64(1000000)
	fmt.Printf("CPU time(usr+sys): %f Sec\n", cpuTime)
}

// GetTotalMemoryUsage runs in a goroutine during the lifetime of a process
func GetTotalMemoryUsage(cmd *exec.Cmd, memLimit uint32, memChan chan<- uint32, wg *sync.WaitGroup) {
	pid := cmd.Process.Pid
	fmt.Println("DEBUG:: Pid to Rip > ", pid)

	defer wg.Done()
	defer close(memChan)

	var wpid int = 0
	var memUsage uint32 = 64 // lower limit for memory used

	var status unix.WaitStatus

	for wpid == 0 {
		fmt.Println("DDDDD==========")
		time.Sleep(1 * time.Millisecond)

		m, err := getMemoryUsage(pid)
		memUsage = maxInt32(memUsage, m)

		if err != nil {
			fmt.Printf("Error in getMemoryUsage: %v\n", err)
			break
		}
		// if memUsage > memLimit {
		// 	err := syscall.Kill(pid, syscall.SIGKILL) // TODO: handle error
		// 	if err != nil {
		// 		fmt.Printf("Error in getMemoryUsage / Limit exceeded: %v\n", err)
		// 	}
		// 	break
		// }
		var werr error
		// wait fro the process to change state
		fmt.Printf("Waiting for pid %d to change state\n", pid)
		wpid, werr = unix.Wait4(pid, &status, (unix.WUNTRACED | unix.WCONTINUED), nil)
		fmt.Printf("DEBUG: Wait4 -> pid: %d , status: %v err: %v\n", wpid, status, werr)
		if werr != nil {
			fmt.Printf("Error in Wait4 call : %v\n", werr)
			break
		}
	}
	memChan <- memUsage
	fmt.Println("USAGE: ", memUsage)
}

func getMemoryUsage(pid int) (uint32, error) {
	fd, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		fmt.Printf("Error in opening: %v", err)
		return 0, err
	}
	defer fd.Close()

	var data, stack uint32 = 0, 0

	pfxData := []byte("VmData")
	pfxStack := []byte("VmStk")
	scanner := bufio.NewScanner(fd)
	fmt.Println("Before scan")
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, pfxData) {
			_, err := fmt.Sscanf(string(line[7:]), "%d", &data)
			if err != nil {
				return 0, err
			}
		}
		if bytes.HasPrefix(line, pfxStack) {
			//get VmStk
			_, err := fmt.Sscanf(string(line[6:]), "%d", &stack)
			if err != nil {
				return 0, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	fmt.Printf("DEBUG: getMemoryUsage -> %d kB\n", data+stack)

	return data + stack, nil
}

func maxInt32(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
