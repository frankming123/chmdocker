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
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0, // 映射为root
				HostID:      syscall.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0, // 映射为root
				HostID:      syscall.Getgid(),
				Size:        1,
			},
		},
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// 传递管道
	cmd.ExtraFiles = []*os.File{rpipe}

	// 设置环境变量
	cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}

	// 组合镜像
	lowerdir := "/var/lib/chmdocker/lower"
	upperdir := "/var/lib/chmdocker/upper"
	workdir := "/var/lib/chmdocker/work"
	merged := "/var/lib/chmdocker/merged"
	o := NewOverlay2([]string{lowerdir}, upperdir, workdir, merged)
	if err := o.Set(); err != nil {
		log.Fatalf("Mount overlay2 error: %v", err)
	}
	//newWorkSpace(rootURL,mntURL)
	// 修改工作目录
	cmd.Dir = merged

	return cmd, wpipe
}
