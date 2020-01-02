package main

import (
	"chmdocker/container"
	log "github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, command string) {
    // 配置新进程环境
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}