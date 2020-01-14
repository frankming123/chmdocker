package cgroups

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
)

// Cgroup 包括路径，资源等信息
type Cgroup struct {
	// Mounts 存放所有挂载点
	Mounts map[string]string

	// Paths 存放所有cgroup路径
	Paths map[string]string

	// ScopePrefix 表示容器域
	ScopePrefix string

	// Resources 包含cgroup资源
	*Resources
}

// NewCgroup 初始化Cgroup结构体
func NewCgroup(scopePrefix string) *Cgroup {
	c := &Cgroup{
		Mounts:      make(map[string]string),
		Paths:       make(map[string]string),
		ScopePrefix: scopePrefix,
		Resources:   &Resources{},
	}
	// c.GetAllMountpoint()
	c.Mounts = GetAllMountpoint()
	return c
}

// Set 设置cgroup资源
func (c *Cgroup) Set() {
	// cpuset.cpus和cpuset.mems需同时配置才能生效，如果有一项缺少，配置相同即可
	if c.CpusetCpus == "" && c.CpusetMems != "" {
		c.CpusetCpus = c.CpusetMems
	}
	if c.CpusetCpus != "" && c.CpusetMems == "" {
		c.CpusetMems = c.CpusetCpus
	}

	t := reflect.TypeOf(c.Resources).Elem()
	v := reflect.ValueOf(c.Resources).Elem()
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i).Interface().(string)
		if val == "" {
			continue
		}
		file := t.Field(i).Tag.Get("file")
		subsystem := t.Field(i).Tag.Get("subsystem")

		if subsystemPath, ok := c.Mounts[subsystem]; ok {
			var fullPath string
			// 从Paths中获取cgroup路径，如果没有，则调用GetCgroupPath函数获取，并保存至Paths中
			if path, ok := c.Paths[subsystem]; ok {
				fullPath = path
			} else {
				fullPath = GetCgroupPath(subsystemPath, c.ScopePrefix, true)
				c.Paths[subsystem] = fullPath
			}
			log.Debugf("Setting cgroup: path: %v file: %v val: %v", fullPath, file, val)

			if err := ioutil.WriteFile(path.Join(fullPath, file), []byte(val), 0644); err != nil {
				log.Errorf("Error set cgroup: write %v fail: %v", fullPath, err)
			}
		} else {
			log.Errorf("Error set cgroup: can not found subsystem path: %v", subsystemPath)
		}
	}
}

// Apply 将pid写入到cgroup的tasks中
func (c *Cgroup) Apply(pid int) {
	for _, fullpath := range c.Paths {
		if err := ioutil.WriteFile(path.Join(fullpath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Errorf("Error set cgroup: write tasks fail: %v", err)
		}
	}
}

// Remove 移除cgroup资源
func (c *Cgroup) Remove() {
	for _, fullpath := range c.Paths {
		log.Debugf("Removing cgroup: path: %v", fullpath)
		if err := os.RemoveAll(fullpath); err != nil {
			log.Errorf("Remove cgroup failed: %v", err)
		}
	}
}

// Resources 包括cgroup资源
type Resources struct {
	// 内存使用量限制
	Memory string `file:"memory.limit_in_bytes" subsystem:"memory"`

	// 设置/读取内存超限控制信息
	OomKillDisable string `file:"memory.oom_control" subsystem:"memory"`

	// 控制各个cgroup组之间的配额占比
	CpuShares string `file:"cpu.shares" subsystem:"cpu"`

	// 限制只能使用特定CPU节点
	CpusetCpus string `file:"cpuset.cpus" subsystem:"cpuset"`

	// 限制只能使用特定内存节点
	CpusetMems string `file:"cpuset.mems" subsystem:"cpuset"`
}
