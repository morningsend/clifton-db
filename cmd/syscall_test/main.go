package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func testStuff1() {
	const SharedMemSize = 1024

	sharedFile, err := os.OpenFile(
		"/tmp/ipc_test_file",
		os.O_CREATE|os.O_RDWR,
		0644,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening file", err)
	}

	err = sharedFile.Truncate(SharedMemSize)

	if err != nil {
		log.Println("error truncate file", err)
	}

	fmt.Println("fd", sharedFile.Fd())

	arena, err := syscall.Mmap(int(sharedFile.Fd()), 0, SharedMemSize, syscall.PROT_WRITE|syscall.PROT_READ,
		syscall.MAP_SHARED)

	if err != nil {
		log.Println("mmap error:", err)
	}

	fmt.Printf("mmap address is %p\n", &arena)

	err = syscall.Munmap(arena)

	pid, err := syscall.ForkExec("/bin/bash", []string{"echo hello world"}, nil)
	fmt.Println("started child process", pid)

	var waitStatus syscall.WaitStatus
	var usage syscall.Rusage

	childPid, err := syscall.Wait4(pid, &waitStatus, 0, &usage)

	fmt.Println("child exited", childPid)

	if waitStatus.Exited() {
		fmt.Println("child process exited with status", waitStatus.ExitStatus())
	}

	if err != nil {
		log.Println(err)
	}
	err = sharedFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func testCmd() {
	cmd := exec.Command("/bin/bash", "-c", "echo 'hello world'")

	out, err := cmd.StdoutPipe()
	checkErr(err)

	errOut, err := cmd.StderrPipe()
	checkErr(err)

	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stderr, errOut)

	err = cmd.Start()
	checkErr(err)
	err = cmd.Wait()
	checkErr(err)
}
func main() {
	testCmd()
}
