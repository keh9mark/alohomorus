//go:build windows
// +build windows

package main

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func setupProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000200, // CREATE_NEW_PROCESS_GROUP
	}
}

func stopProcess(process *os.Process) error {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(process.Pid))
	return cmd.Run()
}