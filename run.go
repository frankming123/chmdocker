package main

import (
	"chmdocker/cgroups"
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, cmd string, res *cgroups.Resources) {
	// 配置新进程环境
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//创建cgroup
	cgroup := cgroups.NewCgroup("chmdocker")
	cgroup.Resources = res
	cgroup.GetAllMountpoint()
	cgroup.Set(parent.Process.Pid)

	parent.Wait()

	log.Debug("container task ended.")
	cgroup.Remove()
	os.Exit(-1)
}
