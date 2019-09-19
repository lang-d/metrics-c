package metrics_c

import "github.com/lang-d/metrics-c/base"

func MemUsePercent() (float32, error) {
	procMeminfo := &base.ProcMeminfo{}
	if err := procMeminfo.Collect(); err != nil {
		return 0, err
	}
	return procMeminfo.UsePercent, nil
}
