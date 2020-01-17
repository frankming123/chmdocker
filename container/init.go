package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// RunContainerInitProcess 子进程的初始化操作，并使用容器进程代替
func RunContainerInitProcess() error {
	cmds := readUserCommand()
	if cmds == nil || len(cmds) == 0 {
		return fmt.Errorf("Run container get user command error, cmd is nil")
	}

	setUpMount()

	path, err := exec.LookPath(cmds[0])
	if err != nil {
		log.Errorf("Exec relative path error: %v", err)
	}
	log.Infof("Find path %s", path)
	// 运行命令
	if err := syscall.Exec(path, cmds[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
