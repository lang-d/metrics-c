package base

import (
	"errors"
	"fmt"
	"github.com/lang-d/metrics-c/g"
	"io/ioutil"
	"runtime"
	"strings"
)

const (
	PROC_STAT_FILE    = "/proc/stat"
	PROC_MEMINFO_FILE = "/proc/meminfo"
	PROC_LOADAVG_FILE = "/proc/loadavg"
	PROC_NET_DEV_FILE = "/proc/net/dev"
)

type CpuMetrics struct {
	Name        string
	User        uint64
	Nice        uint64
	System      uint64
	Idle        uint64
	IoWait      uint64
	Irq         uint64
	Softirq     uint64
	Stealstolen uint64
	Guest       uint64
}

type ProcStat struct {
	Cpus         map[string]CpuMetrics
	CpuIdle      uint64
	Intr         []uint64
	Ctxt         uint64
	Btime        uint64
	Processes    uint64
	ProcsRunning uint64
	ProcsBlocked uint64
	Softirq      []uint64
	TotalCpuTime uint64
}

func (this *ProcStat) parseCpuMetrics(line string) *ProcStat {
	items := strings.Split(line, " ")
	cpu := CpuMetrics{}
	items = g.FilterStringSlice(items, " ")
	for i, item := range items {
		item = strings.TrimSpace(item)
		switch i {
		case 0:
			cpu.Name = item
		case 1:
			cpu.User, _ = g.StringToUint64(item)
		case 2:
			cpu.Nice, _ = g.StringToUint64(item)
		case 3:
			cpu.System, _ = g.StringToUint64(item)
		case 4:
			cpu.Idle, _ = g.StringToUint64(item)
		case 5:
			cpu.IoWait, _ = g.StringToUint64(item)
		case 6:
			cpu.Irq, _ = g.StringToUint64(item)
		case 7:
			cpu.Softirq, _ = g.StringToUint64(item)
		case 8:
			cpu.Stealstolen, _ = g.StringToUint64(item)
		case 9:
			cpu.Guest, _ = g.StringToUint64(item)

		}
	}
	if cpu.Name != "" {
		if this.Cpus != nil {
			this.Cpus[cpu.Name] = cpu
		} else {
			this.Cpus = map[string]CpuMetrics{
				cpu.Name: cpu,
			}
		}

	}
	return this
}

func (this *ProcStat) getMetrics(line string) []uint64 {
	items := strings.Split(line, " ")
	items = g.FilterStringSlice(items, " ")
	metrics := make([]uint64, len(items)-1)
	for i, item := range items {
		item = strings.TrimSpace(item)
		if i > 0 {
			metrics[i-1], _ = g.StringToUint64(item)

		}

	}
	return metrics
}

func (this *ProcStat) getSingleMetric(line string) uint64 {
	items := strings.Split(line, " ")
	items = g.FilterStringSlice(items, " ")
	var metric uint64
	for i, item := range items {
		item = strings.TrimSpace(item)
		if i > 0 {
			metric, _ = g.StringToUint64(item)

		}

	}
	return metric
}

func (this *ProcStat) Collect() error {
	contentBytes, err := ioutil.ReadFile(PROC_STAT_FILE)
	if err != nil {
		return err
	}

	content := g.BytesToString(contentBytes)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu") {
			this.parseCpuMetrics(line)
		}
		if strings.HasPrefix(line, "intr") {
			this.Intr = this.getMetrics(line)
		}
		if strings.HasPrefix(line, "ctxt") {
			this.Ctxt = this.getSingleMetric(line)
		}
		if strings.HasPrefix(line, "btime") {
			this.Btime = this.getSingleMetric(line)
		}
		if strings.HasPrefix(line, "processes") {
			this.Processes = this.getSingleMetric(line)
		}
		if strings.HasPrefix(line, "procs_running") {
			this.ProcsRunning = this.getSingleMetric(line)
		}
		if strings.HasPrefix(line, "procs_blocked") {
			this.ProcsBlocked = this.getSingleMetric(line)
		}
		if strings.HasPrefix(line, "softirq") {
			this.Softirq = this.getMetrics(line)
		}
	}

	if cpu, ok := this.Cpus["cpu"]; ok {
		this.TotalCpuTime += cpu.Softirq + cpu.Guest + cpu.Stealstolen + cpu.Irq + cpu.IoWait + cpu.Idle + cpu.System + cpu.Nice + cpu.User
		this.CpuIdle = cpu.Idle
	} else {
		return errors.New("not found cpu metrics from " + PROC_STAT_FILE)
	}

	return nil

}

type MemoryMetric struct {
	Name  string
	Value uint64
}

type ProcMeminfo struct {
	Mems       map[string]MemoryMetric
	TotalSize  uint64
	TotalUsed  uint64
	TotalFree  uint64
	UsePercent float32
}

func (this *ProcMeminfo) Collect() error {
	contentBytes, err := ioutil.ReadFile(PROC_MEMINFO_FILE)
	if err != nil {
		return err
	}

	content := g.BytesToString(contentBytes)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		items := strings.Fields(line)
		name := strings.Replace(items[0], ":", "", -1)
		memoryMetrics := MemoryMetric{Name: name}
		memoryMetrics.Value, _ = g.StringToUint64(items[1])
		this.Mems[name] = memoryMetrics
	}

	if metric, ok := this.Mems["MemTotal"]; ok {
		this.TotalSize = metric.Value
	} else {
		return errors.New("not found MemToal from " + PROC_MEMINFO_FILE)
	}

	if metric, ok := this.Mems["MemTotal"]; ok {
		this.TotalSize = metric.Value
	} else {
		return errors.New("not found MemToal from " + PROC_MEMINFO_FILE)
	}

	// TotalFree = MemFree + Buffers + Cached
	if metric, ok := this.Mems["MemFree"]; ok {
		this.TotalFree += metric.Value
	} else {
		return errors.New("not found MemFree from " + PROC_MEMINFO_FILE)
	}

	if metric, ok := this.Mems["Buffers"]; ok {
		this.TotalFree += metric.Value
	} else {
		return errors.New("not found Buffers from " + PROC_MEMINFO_FILE)
	}

	if metric, ok := this.Mems["Cached"]; ok {
		this.TotalFree += metric.Value
	} else {
		return errors.New("not found Cached from " + PROC_MEMINFO_FILE)
	}

	this.TotalUsed = this.TotalSize - this.TotalFree
	this.UsePercent = float32(this.TotalUsed) / float32(this.TotalSize)

	return nil
}

type ProcLoadavg struct {
	Load1Min          float32
	Load5Min          float32
	Load15Min         float32
	RunningProcessNum uint64
	TotalProcessNum   uint64
	RunningPid        uint64
	CpuNum            uint8
	Load1             float32
	Load5             float32
	Load15            float32
}

func (this *ProcLoadavg) Collect() error {
	contentBytes, err := ioutil.ReadFile(PROC_LOADAVG_FILE)
	if err != nil {
		return err
	}

	content := g.BytesToString(contentBytes)

	items := strings.Fields(content)
	this.Load1Min, _ = g.StringToFloat32(items[0])
	this.Load5Min, _ = g.StringToFloat32(items[1])
	this.Load15Min, _ = g.StringToFloat32(items[2])
	process := strings.Split(items[3], "/")
	this.RunningProcessNum, _ = g.StringToUint64(process[0])
	this.TotalProcessNum, _ = g.StringToUint64(process[1])
	this.RunningPid, _ = g.StringToUint64(items[4])
	this.CpuNum = uint8(runtime.NumCPU())

	this.Load1 = float32(this.Load1Min) / float32(this.CpuNum)
	this.Load5 = float32(this.Load5Min) / float32(this.CpuNum)
	this.Load15 = float32(this.Load15Min) / float32(this.CpuNum)

	return nil
}

type NetDevMetrics struct {
	Name string

	ReceiveBytes      uint64
	ReceivePackets    uint64
	ReceiveErrs       uint64
	ReceiveDrop       uint64
	ReceiveFifo       uint64
	ReceiveFrame      uint64
	ReceiveCompressed uint64
	ReceiveMulticast  uint64

	TransmitBytes      uint64
	TransmitPackets    uint64
	TransmitErrs       uint64
	TransmitDrop       uint64
	TransmitFifo       uint64
	TransmitColls      uint64
	TransmitCarrier    uint64
	TransmitCompressed uint64
}

type ProcNetDev struct {
	Devs map[string]NetDevMetrics

	TotalReceiveBytes      uint64
	TotalReceivePackets    uint64
	TotalReceiveErrs       uint64
	TotalReceiveDrop       uint64
	TotalReceiveFifo       uint64
	TotalReceiveFrame      uint64
	TotalReceiveCompressed uint64
	TotalReceiveMulticast  uint64

	TotalTransmitBytes      uint64
	TotalTransmitPackets    uint64
	TotalTransmitErrs       uint64
	TotalTransmitDrop       uint64
	TotalTransmitFifo       uint64
	TotalTransmitColls      uint64
	TotalTransmitCarrier    uint64
	TotalTransmitCompressed uint64
}

func (this *ProcNetDev) Collect() error {
	contentBytes, err := ioutil.ReadFile(PROC_NET_DEV_FILE)
	if err != nil {
		return err
	}

	content := g.BytesToString(contentBytes)

	lines := strings.Split(content, "\n")
	lines = lines[2:]
	for _, line := range lines {
		items := strings.Fields(line)
		netDevMetrics := NetDevMetrics{}
		for i, item := range items {
			switch i {
			case 0:
				netDevMetrics.Name = strings.Replace(item, ":", "", -1)
			case 1:
				netDevMetrics.ReceiveBytes, _ = g.StringToUint64(item)
				this.TotalReceiveBytes += netDevMetrics.ReceiveBytes
			case 2:
				netDevMetrics.ReceivePackets, _ = g.StringToUint64(item)
				this.TotalReceivePackets += netDevMetrics.ReceivePackets
			case 3:
				netDevMetrics.ReceiveErrs, _ = g.StringToUint64(item)
				this.TotalReceiveErrs += netDevMetrics.ReceiveErrs
			case 4:
				netDevMetrics.ReceiveDrop, _ = g.StringToUint64(item)
				this.TotalReceiveDrop += netDevMetrics.ReceiveDrop
			case 5:
				netDevMetrics.ReceiveFifo, _ = g.StringToUint64(item)
				this.TotalReceiveFifo += netDevMetrics.ReceiveFifo
			case 6:
				netDevMetrics.ReceiveFrame, _ = g.StringToUint64(item)
				this.TotalReceiveFrame += netDevMetrics.ReceiveFrame
			case 7:
				netDevMetrics.ReceiveCompressed, _ = g.StringToUint64(item)
				this.TotalReceiveCompressed += netDevMetrics.ReceiveCompressed
			case 8:
				netDevMetrics.ReceiveMulticast, _ = g.StringToUint64(item)
				this.TotalReceiveMulticast += netDevMetrics.ReceiveMulticast
			case 9:
				netDevMetrics.TransmitBytes, _ = g.StringToUint64(item)
				this.TotalTransmitBytes += netDevMetrics.TransmitBytes
			case 10:
				netDevMetrics.TransmitPackets, _ = g.StringToUint64(item)
				this.TotalTransmitPackets += netDevMetrics.TransmitPackets
			case 11:
				netDevMetrics.TransmitErrs, _ = g.StringToUint64(item)
				this.TotalTransmitErrs += netDevMetrics.TransmitErrs
			case 12:
				netDevMetrics.TransmitDrop, _ = g.StringToUint64(item)
				this.TotalTransmitDrop += netDevMetrics.TransmitDrop
			case 13:
				netDevMetrics.TransmitFifo, _ = g.StringToUint64(item)
				this.TotalTransmitFifo += netDevMetrics.TransmitFifo
			case 14:
				netDevMetrics.TransmitColls, _ = g.StringToUint64(item)
				this.TotalTransmitColls += netDevMetrics.TransmitColls
			case 15:
				netDevMetrics.TransmitCarrier, _ = g.StringToUint64(item)
				this.TotalTransmitCarrier += netDevMetrics.TransmitCarrier
			case 16:
				netDevMetrics.TransmitCompressed, _ = g.StringToUint64(item)
				this.TotalTransmitCompressed += netDevMetrics.TransmitCompressed
			}
		}
		this.Devs[netDevMetrics.Name] = netDevMetrics
	}

	return nil
}

type ProcPidFd struct {
	Pid   uint64
	FdNum uint64
}

func (this *ProcPidFd) Collect() error {
	if this.Pid == 0 {
		return errors.New("/proc/pid/fd dir need a pid,but not give")
	}

	fdPath := fmt.Sprintf("/prc/%d/fd", this.Pid)

	files, err := ioutil.ReadDir(fdPath)
	if err != nil {
		return err
	}
	this.FdNum = uint64(len(files))

	return nil
}
