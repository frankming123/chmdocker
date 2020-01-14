package main

import (
	"chmdocker/cgroups"
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

// Run run子命令细节
func Run(tty bool, cmd string, res *cgroups.Resources) {
	// 配置新进程环境
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// cpuset.cpus和cpuset.mems需同时配置才能生效，如果有一项缺少，配置相同即可
	if res.CpusetCpus == "" && res.CpusetMems != "" {
		res.CpusetCpus = res.CpusetMems
	}
	if res.CpusetCpus != "" && res.CpusetMems == "" {
		res.CpusetMems = res.CpusetCpus
	}

	//创建cgroup
	cgroup := cgroups.NewCgroup("chmdocker")
	cgroup.Resources = res
	cgroup.Set()
	cgroup.Apply(parent.Process.Pid)

	parent.Wait()

	log.Debug("container task ended.")
	cgroup.Remove()
	os.Exit(-1)
}
