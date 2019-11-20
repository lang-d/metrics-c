package metrics_c

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_CPU_MONITOR_TIME = time.Second * 3

func CpuUsePercent() (float32, error) {
	return CalcCpuUsePercent(DEFAULT_CPU_MONITOR_TIME)
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func CalcCpuUsePercent(sleepTime time.Duration) (float32, error) {
	//procStat1 := &base.ProcStat{}
	//if err := procStat1.Collect(); err != nil {
	//	return 0, err
	//}
	//time.Sleep(sleepTime)
	//procStat2 := &base.ProcStat{}
	//if err := procStat2.Collect(); err != nil {
	//	return 0, err
	//}
	//
	//// pcpu =100* idle/total
	//total := procStat2.TotalCpuTime - procStat1.TotalCpuTime
	//idle := procStat2.CpuIdle - procStat1.CpuIdle
	//
	//if total == 0 {
	//	return 0, nil
	//}
	//
	//return 100 * float32(idle) / float32(total), nil
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()
	idleTicks := float32(idle1 - idle0)
	totalTicks := float32(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
	return cpuUsage, nil
}
