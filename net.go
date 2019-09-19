package metrics_c

import (
	"github.com/lang-d/metrics-c/base"
	"time"
)

type NetMetrics struct {
	ReceiveSpeed float32
	SendSpeed    float32
	SpeedUnit    string
}

const DEFAULT_NET_MONITOR_TIME = time.Second * 5

func NetStatus() (*NetMetrics, error) {
	return CalcNetStaus(DEFAULT_NET_MONITOR_TIME)
}

func CalcNetStaus(sleepTime time.Duration) (*NetMetrics, error) {
	procNetDev1 := &base.ProcNetDev{}
	if err := procNetDev1.Collect(); err != nil {
		return nil, err
	}

	time.Sleep(sleepTime)

	procNetDev2 := &base.ProcNetDev{}
	if err := procNetDev2.Collect(); err != nil {
		return nil, err
	}

	received := (procNetDev2.TotalReceiveBytes - procNetDev1.TotalReceiveBytes) / 1024
	sended := (procNetDev2.TotalTransmitBytes - procNetDev1.TotalTransmitBytes) / 1024

	receiveSpeed := float32(received) / float32(sleepTime.Seconds())
	sendSpeed := float32(sended) / float32(sleepTime.Seconds())
	return &NetMetrics{
		ReceiveSpeed: receiveSpeed,
		SendSpeed:    sendSpeed,
		SpeedUnit:    "kb/s",
	}, nil

}
