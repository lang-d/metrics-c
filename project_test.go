package metrics_c_test

import (
	"fmt"
	"github.com/lang-d/metrics-c"
	"testing"
)

const pid = 1151

func TestCpuUsePercent(t *testing.T) {
	cpuUse, err := metrics_c.CpuUsePercent()
	if err != nil {
		panic(err)
	}
	fmt.Printf("cpu use %f", cpuUse)
}

func TestDiskUsePercent(t *testing.T) {
	diskUse, err := metrics_c.DiskUsePercent()
	if err != nil {
		panic(err)
	}
	fmt.Printf("disk use %f", diskUse)
}

func TestMemUsePercent(t *testing.T) {
	memUse, err := metrics_c.MemUsePercent()
	if err != nil {
		panic(err)
	}
	fmt.Printf("mem use %f", memUse)
}

func TestNetStatus(t *testing.T) {
	net, err := metrics_c.NetStatus()
	if err != nil {
		panic(err)
	}
	fmt.Printf("net %v", net)
}

func TestLoadAvg(t *testing.T) {
	load, err := metrics_c.LoadAvg()
	if err != nil {
		panic(err)
	}
	fmt.Printf("load %v", load)
}

func TestProcessAll(t *testing.T) {
	process, err := metrics_c.ProcessAll(pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("process %v", process)
}

func TestProcessFdNum(t *testing.T) {
	fdNum, err := metrics_c.ProcessFdNum(pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("process %d fd num %d", pid, fdNum)
}

func TestProcessCpuUsePercent(t *testing.T) {
	cpuUse, err := metrics_c.ProcessCpuUsePercent(pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("process %d cpu use %f ", pid, cpuUse)
}

func TestProcessMemUsePercent(t *testing.T) {
	memUse, err := metrics_c.ProcessMemUsePercent(pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("process %d mem use %f", pid, memUse)
}
