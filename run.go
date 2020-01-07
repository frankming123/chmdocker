package main

import (
	"chmdocker/cgroups"
	"chmdocker/cgroups/subsystems"
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, cmd string, res *subsystems.ResourceConfig) {
	// 配置新进程环境
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//创建cgroup
	cgroupManager := cgroups.NewCgroupManager("chmdocker-cgroup")
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	parent.Wait()

	log.Debug("container task ended.")
	cgroupManager.Destroy()
	os.Exit(-1)
}
