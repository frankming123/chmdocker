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

	// 挂载rootfs，指定为独立的mount命名空间（默认为共享）
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// 挂载/proc目录
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

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
