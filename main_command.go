package main

import (
	"chmdocker/cgroups"
	"chmdocker/container"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "memory, m",
			Usage: "Memory limit (format: <number>[<unit>], where unit = b, k, m or g)",
		},
		cli.StringFlag{
			Name:  "cpu-shares",
			Usage: "CPU shares",
		},
		cli.StringFlag{
			Name:  "cpuset-cpus",
			Usage: "CPUs in which to allow execution (0-3, 0,1)",
		},
		cli.StringFlag{
			Name:  "cpuset-mems",
			Usage: "Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on NUMA systems.",
		},
	},
	// 输入参数后执行的操作
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := c.Args().Get(0)
		tty := c.Bool("it")
		resource := &cgroups.Resources{
			Memory:     c.String("memory"),
			CpuShares:   c.String("cpu-shares"),
			CpusetCpus: c.String("cpuset-cpus"),
			CpusetMems: c.String("cpuset-mems"),
		}
		// 运行Run函数
		Run(tty, cmd, resource)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(c *cli.Context) error {
		log.Debugf("init come on")
		cmd := c.Args().Get(0)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
