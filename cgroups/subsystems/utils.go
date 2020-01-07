package subsystems

import (
	"bufio"
	"fmt"
	_ "github.com/Sirupsen/logrus"
	"os"
	"path"
	"strings"
)

// FindCgroupMountpoint 查找cgroup指定子系统的挂载路径
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// GetCgroupPath 拼接将要创建的cgroup绝对路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)

	fullPath := path.Join(cgroupRoot, cgroupPath)

	if _, err := os.Stat(fullPath); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(fullPath, 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}

		}
		return fullPath, nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
