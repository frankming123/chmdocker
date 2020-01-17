package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 创建子进程，并隔离命名空间
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	rpipe, wpipe, err := os.Pipe()
	if err != nil {
		log.Errorf("New pipe error: %v", err)
		return nil, nil
	}
	// /proc/self/是一个链接，指向进程自身，/proc/PID/
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// 传递管道
	cmd.ExtraFiles = []*os.File{rpipe}

	// 修改工作目录
	cmd.Dir = "/tmp/alpine"

	return cmd, wpipe
}
