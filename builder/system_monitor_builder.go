package builder

import (
	"monitor-system/monitoring"
	"time"
)

type SystemMonitorBuilder struct {
	Interval time.Duration
}

func NewSystemMonitorBuilder() *SystemMonitorBuilder {
	return &SystemMonitorBuilder{
		Interval: time.Minute,
	}
}

func (builder *SystemMonitorBuilder) WithInterval(interval time.Duration) *SystemMonitorBuilder {
	builder.Interval = interval
	return builder
}

func (builder *SystemMonitorBuilder) Build() *monitoring.SystemMonitor {
	return monitoring.NewSystemMonitor(builder.Interval)
}
