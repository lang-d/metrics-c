package base

import (
	"errors"
	"fmt"
	"github.com/lang-d/metrics-c/g"
	"strings"
)

const (
	DF_CMD = "df"
)

type DiskMetrics struct {
	Filesystem string
	Size       uint64
	Used       uint64
	Available  uint64
	Use        uint64
	MountedOn  string
}

type Df struct {
	Filesystems    map[string]DiskMetrics
	TotalSize      uint64
	TotalUsed      uint64
	TotalAvailable uint64
	UsePercent     float32
}

func (this *Df) Collect() error {
	content, err := g.DoCmdAndOutPut(DF_CMD)
	if err != nil {
		return err
	}

	lines := strings.Split(content, "\n")
	this.Filesystems = map[string]DiskMetrics{}
	for i, line := range lines {
		if i > 0 {
			items := strings.Fields(line)
			diskMetrics := DiskMetrics{}
			for k, item := range items {
				switch k {
				case 0:
					diskMetrics.Filesystem = item
				case 1:
					diskMetrics.Size, _ = g.StringToUint64(item)
					this.TotalSize += diskMetrics.Size
				case 2:
					diskMetrics.Used, _ = g.StringToUint64(item)
					this.TotalUsed += diskMetrics.Used
				case 3:
					diskMetrics.Available, _ = g.StringToUint64(item)
					this.TotalAvailable += diskMetrics.Available
				case 4:
					diskMetrics.Use, _ = g.StringToUint64(strings.Replace(item, "%", "", -1))
				case 5:
					diskMetrics.MountedOn = item

				}
			}
			this.Filesystems[diskMetrics.Filesystem] = diskMetrics
		}
	}

	this.UsePercent = float32(this.TotalUsed) / float32(this.TotalSize)

	return nil
}

// ps -p pid -u
type PsPPidU struct {
	User          string
	Pid           uint64 // Pid must set before collect
	CpuUsePercent float32
	MemUsePercent float32
	Vsz           uint64
	Rss           uint64
	TTy           string
	Stat          string
	Start         string
	Time          string
	Command       string
}

func (this *PsPPidU) Collect() error {
	if this.Pid == 0 {
		return errors.New("ps -p pid -u command need a pid,but not give")
	}

	cmd := fmt.Sprintf("ps -p %d -u", this.Pid)
	outPut, err := g.DoCmdAndOutPut(cmd)
	if err != nil {
		return err
	}
	lines := strings.Split(outPut, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("not found process info by pid %d", this.Pid)
	}

	items := strings.Fields(lines[1])
	for i, item := range items {
		// USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
		switch i {
		case 0:
			this.User = item
		case 2:
			this.CpuUsePercent, _ = g.StringToFloat32(item)
		case 3:
			this.MemUsePercent, _ = g.StringToFloat32(item)
		case 4:
			this.Vsz, _ = g.StringToUint64(item)
		case 5:
			this.Rss, _ = g.StringToUint64(item)
		case 6:
			this.TTy = item
		case 7:
			this.Stat = item
		case 8:
			this.Start = item
		case 9:
			this.Time = item
		case 10:
			this.Command = item

		}
	}

	return nil

}
