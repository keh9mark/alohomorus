//go:build !windows
// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func setupProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func stopProcess(process *os.Process) error {
	return syscall.Kill(-process.Pid, syscall.SIGKILL)
}