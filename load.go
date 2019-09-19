package metrics_c

import "github.com/lang-d/metrics-c/base"

type LoadAvgMetrics struct {
	Load1Min  float32
	Load5Min  float32
	Load15Min float32
	Load1     float32
	Load5     float32
	Load15    float32
	CpuNum    uint8
}

func LoadAvg() (*LoadAvgMetrics, error) {
	procLoad := &base.ProcLoadavg{}
	if err := procLoad.Collect(); err != nil {
		return nil, err
	}

	return &LoadAvgMetrics{
		Load1Min:  procLoad.Load1Min,
		Load5Min:  procLoad.Load5Min,
		Load15Min: procLoad.Load15Min,
		Load1:     procLoad.Load1,
		Load5:     procLoad.Load5,
		Load15:    procLoad.Load15,
		CpuNum:    procLoad.CpuNum,
	}, nil

}
