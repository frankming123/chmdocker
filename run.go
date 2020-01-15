package main

import (
	"chmdocker/cgroups"
	"strings"
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

// Run run子命令细节
func Run(tty bool, cmds []string, res *cgroups.Resources) {
	// 配置新进程环境
	parent, wpipe := container.NewParentProcess(tty)
	if parent==nil{
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//创建cgroup
	cgroup := cgroups.NewCgroup("chmdocker")
	cgroup.Resources = res
	cgroup.Set()
	cgroup.Apply(parent.Process.Pid)

	// 发送用户命令
	sendInitCommand(cmds,wpipe)

	parent.Wait()

	log.Debug("container task ended.")
	cgroup.Remove()
	os.Exit(-1)
}

func sendInitCommand(cmds []string,wpipe *os.File){
	cmd:=strings.Join(cmds," ")
	log.Infof("container command is %s",cmd)
	wpipe.WriteString(cmd)
	wpipe.Close()
}