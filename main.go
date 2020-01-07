package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `chmdocker is a simple container runtime implementation.
			   Enjoy it, just for fun.`

func main() {
	// 初始化命令行参数的配置
	app := cli.NewApp()
	app.Name = "chmdocker"
	app.Version = "0.0.1"
	app.Usage = usage

	// 配置子命令：init，run
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	// 运行app之前的操作，这里主要是配置log
	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)
		return nil
	}

	// 运行app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
