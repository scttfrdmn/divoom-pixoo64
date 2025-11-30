package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemMetrics struct {
	CPUPercent    float64
	MemoryPercent float64
	MemoryUsedGB  float64
	MemoryTotalGB float64
	NetSentMB     float64
	NetRecvMB     float64
	Timestamp     time.Time
}

type Collector struct {
	lastNetStat   *net.IOCountersStat
	lastCheckTime time.Time
}

func NewCollector() *Collector {
	return &Collector{
		lastCheckTime: time.Now(),
	}
}

func (c *Collector) Collect() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// Get CPU usage
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("get cpu percent: %w", err)
	}
	if len(cpuPercents) > 0 {
		metrics.CPUPercent = cpuPercents[0]
	}

	// Get memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get memory stats: %w", err)
	}
	metrics.MemoryPercent = vmStat.UsedPercent
	metrics.MemoryUsedGB = float64(vmStat.Used) / 1024 / 1024 / 1024
	metrics.MemoryTotalGB = float64(vmStat.Total) / 1024 / 1024 / 1024

	// Get network stats
	netStats, err := net.IOCounters(false)
	if err == nil && len(netStats) > 0 {
		currentStat := &netStats[0]

		if c.lastNetStat != nil {
			// Calculate rate since last check
			timeDiff := time.Since(c.lastCheckTime).Seconds()
			if timeDiff > 0 {
				bytesSent := float64(currentStat.BytesSent - c.lastNetStat.BytesSent)
				bytesRecv := float64(currentStat.BytesRecv - c.lastNetStat.BytesRecv)

				metrics.NetSentMB = (bytesSent / 1024 / 1024) / timeDiff
				metrics.NetRecvMB = (bytesRecv / 1024 / 1024) / timeDiff
			}
		}

		c.lastNetStat = currentStat
		c.lastCheckTime = time.Now()
	}

	return metrics, nil
}

func (m *SystemMetrics) String() string {
	return fmt.Sprintf(
		"CPU: %.1f%% | MEM: %.1f%% (%.1f/%.1fGB) | NET: ↓%.2fMB/s ↑%.2fMB/s",
		m.CPUPercent,
		m.MemoryPercent,
		m.MemoryUsedGB,
		m.MemoryTotalGB,
		m.NetRecvMB,
		m.NetSentMB,
	)
}
