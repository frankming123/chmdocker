package container

import (
	log "github.com/Sirupsen/logrus"
	"strings"
	"io/ioutil"
	"os"
)

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	return strings.Split(string(msg)," ")
}
