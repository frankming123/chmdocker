package container

import (
	"github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

    // 挂载rootfs，指定为独立的mount命名空间（默认为共享）
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    // 挂载/proc目录
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
    argv := []string{command}
    // 运行命令
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}