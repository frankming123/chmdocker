package cgroups

import (
	"bufio"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// GetAllMountpoint 查找cgroup指定子系统的挂载路径，并写入到Mounts中
func GetAllMountpoint() (mounts map[string]string) {
	mounts = make(map[string]string)
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Errorf("Get all mountpoint error: %v", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if _, ok := mounts[opt]; !ok {
				mounts[opt] = fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Warnf("Get all mountpoint error: mount point scanner error: %v", err)
	}
	return
}

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

// copyCpuOrMemIfNeeded 从父目录复制cpuset.cpus和cpuset.mems，如果当前目录的对应文件为空
func copyCpuOrMemIfNeeded(current string) {
	parent := path.Dir(current)
	ccpu, _ := ioutil.ReadFile(path.Join(current, "cpuset.cpus"))
	if isEmpty(ccpu) {
		pcpu, _ := ioutil.ReadFile(path.Join(parent, "cpuset.cpus"))
		if err := ioutil.WriteFile(path.Join(current, "cpuset.cpus"), pcpu, 0644); err != nil {
			log.Errorf("Error copy cpuset.cpus: write %v fail: %v", current, err)
		}
	}
	cmem, _ := ioutil.ReadFile(path.Join(current, "cpuset.mems"))
	if isEmpty(cmem) {
		pmem, _ := ioutil.ReadFile(path.Join(parent, "cpuset.mems"))
		if err := ioutil.WriteFile(path.Join(current, "cpuset.mems"), pmem, 0644); err != nil {
			log.Errorf("Error copy cpuset.mems: write %v fail: %v", current, err)
		}
	}

}

func isEmpty(b []byte) bool {
	return len(bytes.Trim(b, "\n")) == 0
}
