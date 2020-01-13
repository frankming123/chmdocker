package cgroups

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

// Cgroup 包括路径，资源等信息
type Cgroup struct {
	// Paths 存放各个subsystem的绝对路径
	Paths map[string]string

	// ScopePrefix 表示容器域
	ScopePrefix string

	// Resources 包含cgroup资源
	*Resources
}

// NewCgroup 初始化Cgroup结构体
func NewCgroup(scopePrefix string) *Cgroup {
	c := &Cgroup{
		Paths:       make(map[string]string),
		ScopePrefix: scopePrefix,
		Resources:   &Resources{},
	}
	return c
}

// GetAllMountpoint 查找cgroup指定子系统的挂载路径，并写入到Paths中
func (c *Cgroup) GetAllMountpoint() {
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
			if _, ok := c.Paths[opt]; !ok {
				c.Paths[opt] = fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Warnf("Get all mountpoint error: mount point scanner error: %v", err)
	}
}

// Set 设置cgroup资源
func (c *Cgroup) Set(pid int) {
	t := reflect.TypeOf(c.Resources).Elem()
	v := reflect.ValueOf(c.Resources).Elem()
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i).Interface().(string)
		if val == "" {
			continue
		}
		file := t.Field(i).Tag.Get("file")
		subsystem := t.Field(i).Tag.Get("subsystem")

		if subsystemPath, ok := c.Paths[subsystem]; ok {
			fullPath := GetCgroupPath(subsystemPath, c.ScopePrefix, true)
			log.Debugf("Setting cgroup: subsystem: %v file: %v val: %v", fullPath, file, val)

			if err := ioutil.WriteFile(path.Join(fullPath, file), []byte(val), 0644); err != nil {
				log.Errorf("Error set cgroup: write %v fail: %v", fullPath, err)
			}
			if err := ioutil.WriteFile(path.Join(fullPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
				log.Errorf("Error set cgroup: write tasks fail: %v", err)
			}
		} else {
			log.Errorf("Error set cgroup: can not found subsystem path: %v", subsystemPath)
		}
	}
}

// Remove 移除cgroup资源
func (c *Cgroup) Remove() {
	t := reflect.TypeOf(c.Resources).Elem()
	v := reflect.ValueOf(c.Resources).Elem()
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i).Interface().(string)
		if val == "" {
			continue
		}
		subsystem := t.Field(i).Tag.Get("subsystem")

		if subsystemPath, ok := c.Paths[subsystem]; ok {
			fullPath := GetCgroupPath(subsystemPath, c.ScopePrefix, true)
			log.Debugf("Removing cgroup: subsystem: %v", fullPath)
			if err := os.RemoveAll(fullPath); err != nil {
				log.Errorf("Remove cgroup failed: %v", err)
			}
		}
	}
}

// Resources 包括cgroup资源
type Resources struct {
	// Memory limit (in bytes)
	Memory string `file:"memory.limit_in_bytes" subsystem:"memory"`

	// Whether to disable OOM Killer
	OomKillDisable string `file:"memory.oom_control" subsystem:"memory"`

	// CPU shares (relative weight vs. other containers)
	CpuShares string `file:"cpu.shares" subsystem:"cpu"`

	// CPU to use
	CpusetCpus string `file:"cpuset.cpus" subsystem:"cpuset"`
}
