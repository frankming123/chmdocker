package cgroups

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

// Subsystem 子系统名称
type Subsystem int

const (
	// BLKIO subsystem: 为块设备设定输入/输出限制,比如物理驱动设备(包括磁盘,固态硬盘,USB等)
	BLKIO Subsystem = iota
	// CPU subsystem: 使用调度程序控制task对CPU的使用
	CPU
	// CPUACCT subsystem: 自动生成cgroup中task对CPU资源使用情况的报告
	CPUACCT
	// CPUSET subsystem: 为cgroup中的task分配独立的CPU(此处针对多处理器系统)和内存
	CPUSET
	// DEVICES subsystem: 开启或关闭cgroup中task对设备的访问
	DEVICES
	// FREEZER subsystem: 挂起或恢复cgroup中的task
	FREEZER
	// MEMORY subsystem: 设定cgroup中task对内存使用量的限定,并且自动生成这些task对内存资源使用情况的报告
	MEMORY
	// PERFEVENT subsystem: 使得cgroup中的task可以进行统一的性能测试
	PERFEVENT
	// NETCLS subsystem: 通过使用等级识别符(classid)标记网络数据包,从而允许Linux流量控制程序(TC：Traffic Controller)识别从具体cgroup中生成的数据包
	NETCLS
)

// String 将Subsystem转换为string类型
func (s Subsystem) String() string {
	switch s {
	case BLKIO:
		return "blkio"
	case CPU:
		return "cpu"
	case CPUACCT:
		return "cpuacct"
	case CPUSET:
		return "cpuset"
	case DEVICES:
		return "devices"
	case FREEZER:
		return "freezer"
	case MEMORY:
		return "memory"
	case PERFEVENT:
		return "perfevent"
	case NETCLS:
		return "net_cls"
	}
	return "unknown"
}

// FindCgroupMountpoint 查找cgroup指定子系统的挂载路径
func (s Subsystem) FindCgroupMountpoint() string {
	if s.String() == "unknown" {
		log.Warnln("Can not find cgroup mount point: subsystem \"unknown\" is not exist.")
		return ""
	}

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
			if opt == s.String() {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Warnf("Can not find cgroup mount point: mount point scanner error: %v", err)
		return ""
	}

	log.Warnf("Can not find cgroup mount point: subsystem %v not found", s.String())
	return ""
}
