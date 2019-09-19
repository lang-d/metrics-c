package metrics_c

import "github.com/lang-d/metrics-c/base"

func DiskUsePercent() (float32, error) {
	df := &base.Df{}
	if err := df.Collect(); err != nil {
		return 0, nil
	}
	return df.UsePercent, nil
}
