package main

import "flag"

// Config struct for flags
type Config struct {
	CPU    int
	AS     int
	Memory int
	Nproc  int
	Fsize  int
	Stack  int
	Clock  int
	Chroot string
	Exec   string
}

// HandleFlags handles all the flags passed to the program
func HandleFlags() *Config {
	conf := new(Config) // create a new pointer to Config struct

	flag.IntVar(&conf.CPU, "cpu", 5, "CPU time in seconds")
	flag.IntVar(&conf.AS, "addrSpcae", 65536, "Adress Space Memory in KBytes")
	flag.IntVar(&conf.Memory, "mem", 32768, "Program Memory in KBytes")
	flag.IntVar(&conf.Nproc, "nproc", 10, "Number of processes the program can create")
	flag.IntVar(&conf.Fsize, "fsize", 8192, "Output file size in Kbytes")
	flag.IntVar(&conf.Stack, "stack", 8192, "Stack size in Kbytes")
	flag.IntVar(&conf.Clock, "clock", 10, "Wall Clock time limit in seconds")
	flag.StringVar(&conf.Chroot, "chroot", "/", "Chroot folder path. Should have a valid root FileSystem")
	flag.StringVar(&conf.Exec, "exec", "", "Executable path")
	flag.Parse()

	// TODO: Check String Flags for errors
	conf.Memory *= 1024 // convert Bytes to Kbytes
	conf.AS *= 1024
	conf.Fsize *= 1024
	conf.Stack *= 1024

	return conf
}
