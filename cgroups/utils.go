package cgroups

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
)

// GetCgroupPath 拼接将要创建的cgroup绝对路径
func GetCgroupPath(p string, cgroupPath string, autoCreate bool) string {
	cgroupRoot := p

	fullPath := path.Join(cgroupRoot, cgroupPath)
	_, err := os.Stat(fullPath)
	if err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(fullPath, 0755); err != nil {
				log.Errorf("Error create cgroup: %v", err)
				return ""
			}

		}
		return fullPath
	}

	log.Errorf("Error cgroup path: %v", err)
	return ""
}
