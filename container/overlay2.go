package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Overlay2 描述容器overlay2文件系统挂载信息
type Overlay2 struct {
	lowerdir []string
	upperdir string
	merged   string
	workdir  string
}

// NewOverlay2 初始化overlay2结构体
func NewOverlay2(lowerdir []string, upperdir, workdir, merged string) *Overlay2 {
	o := &Overlay2{
		lowerdir: lowerdir,
		upperdir: upperdir,
		merged:   merged,
		workdir:  workdir,
	}
	return o
}

// Set 设置overlay2
func (o *Overlay2) Set() error {
	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", strings.Join(o.lowerdir, ":"), o.upperdir, o.workdir)
	println(data)
	if err := syscall.Mount("overlay", o.merged, "overlay", syscall.MS_RELATIME, data); err != nil {
		return fmt.Errorf("Mount overlay2 failed: %v", err)
	}
	return nil
}

// Remove 移除overlay2
func (o *Overlay2) Remove() error {
	if err := syscall.Unmount(o.merged, 0); err != nil {
		return fmt.Errorf("Unmount overlay2 failed: %v", err)
	}
	return nil
}

// EnsureDirExists 确保目录提前创建
func (o *Overlay2) EnsureDirExists() error {
	args := []string{"-p", o.upperdir, o.workdir, o.merged}
	args = append(args, o.lowerdir...)
	cmd := exec.Command("mkdir", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("EnsureDirExists: Run mkdir failed: %v", err)
	}
	return nil
}
