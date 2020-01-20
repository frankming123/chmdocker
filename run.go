package main

import (
	"chmdocker/cgroups"
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"strings"
	"syscall"
)

// Run run子命令细节
func Run(tty bool, cmds []string, res *cgroups.Resources) {
	// 容器名称
	scope := gensha256()

	containerDir := path.Join("/var/lib/chmdocker/container", scope)
	if err := os.Mkdir(containerDir, 0777); err != nil {
		log.Fatalf("Create containerDir %v failed: %v", containerDir, err)
	}
	// 组合镜像
	lowerdir := path.Join(containerDir, "lower")
	if err := syscall.Symlink("/var/lib/chmdocker/image/alpine", lowerdir); err != nil {
		log.Fatalf("Create symlink failed: %v", err)
	}
	upperdir := path.Join(containerDir, "upper")
	workdir := path.Join(containerDir, "work")
	merged := path.Join(containerDir, "merged")
	overlay2 := container.NewOverlay2([]string{lowerdir}, upperdir, workdir, merged)
	if err := overlay2.Set(); err != nil {
		log.Fatalf("Mount overlay2 error: %v", err)
	}
	// 配置新进程环境
	parent, wpipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	// 修改工作目录
	parent.Dir = merged

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//创建cgroup
	cgroup := cgroups.NewCgroup(scope, res)
	cgroup.Set()
	cgroup.Apply(parent.Process.Pid)

	// 发送用户命令
	sendInitCommand(cmds, wpipe)

	parent.Wait()

	log.Debug("container task ended.")

	cgroup.Remove()
	if err := overlay2.Remove(); err != nil {
		log.Errorf("Remove overlay2 failed: %v", err)
	}
	if err := os.RemoveAll(containerDir); err != nil {
		log.Errorf("Remove containerdir failed: %v", err)
	}
	os.Exit(-1)
}

func sendInitCommand(cmds []string, wpipe *os.File) {
	cmd := strings.Join(cmds, " ")
	log.Infof("container command is %s", cmd)
	wpipe.WriteString(cmd)
	wpipe.Close()
}
