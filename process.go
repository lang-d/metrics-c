package metrics_c

import "github.com/lang-d/metrics-c/base"

func ProcessCpuUsePercent(pid uint64) (float32, error) {
	psPPidU := &base.PsPPidU{Pid: pid}
	if err := psPPidU.Collect(); err != nil {
		return 0, err
	}
	return psPPidU.CpuUsePercent, nil
}

func ProcessMemUsePercent(pid uint64) (float32, error) {
	psPPidU := &base.PsPPidU{Pid: pid}
	if err := psPPidU.Collect(); err != nil {
		return 0, err
	}
	return psPPidU.MemUsePercent, nil
}

func ProcessFdNum(pid uint64) (uint64, error) {
	procPidFd := &base.ProcPidFd{Pid: pid}
	if err := procPidFd.Collect(); err != nil {
		return 0, err
	}
	return procPidFd.FdNum, nil
}

type ProcessMetrics struct {
	Pid           uint64
	Command       string
	MemUsePercent float32
	CpuUsePercent float32
	FdNum         uint64
}

func ProcessAll(pid uint64) (*ProcessMetrics, error) {
	procPidFd := &base.ProcPidFd{Pid: pid}
	psPPidU := &base.PsPPidU{Pid: pid}
	if err := psPPidU.Collect(); err != nil {
		return nil, err
	}
	if err := procPidFd.Collect(); err != nil {
		return nil, err
	}

	return &ProcessMetrics{
		Pid:           pid,
		Command:       psPPidU.Command,
		CpuUsePercent: psPPidU.CpuUsePercent,
		MemUsePercent: psPPidU.MemUsePercent,
		FdNum:         procPidFd.FdNum,
	}, nil
}
