package metrics_c

import (
	"github.com/lang-d/metrics-c/base"
	"time"
)

const DEFAULT_CPU_MONITOR_TIME = time.Second * 3

func CpuUsePercent() (float32, error) {
	return CalcCpuUsePercent(DEFAULT_CPU_MONITOR_TIME)
}

func CalcCpuUsePercent(sleepTime time.Duration) (float32, error) {
	procStat1 := &base.ProcStat{}
	if err := procStat1.Collect(); err != nil {
		return 0, err
	}
	time.Sleep(sleepTime)
	procStat2 := &base.ProcStat{}
	if err := procStat2.Collect(); err != nil {
		return 0, err
	}

	// pcpu =100* idle/total
	total := procStat2.TotalCpuTime - procStat1.TotalCpuTime
	idle := procStat2.CpuIdle - procStat1.CpuIdle

	if total == 0 {
		return 0, nil
	}

	return 100 * float32(idle) / float32(total), nil

}
