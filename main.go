package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
}

func getMetrics() Metrics {
	// Get CPU Usage
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory Usage
	vmStat, _ := mem.VirtualMemory()

	// Get disk usage
	diskUsage, _ := disk.Usage("/")

	return Metrics{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: vmStat.UsedPercent,
		DiskUsage:   diskUsage.UsedPercent,
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := getMetrics()
	tmpl := template.Must(template.ParseFiles("templates/metrics.html"))
	tmpl.Execute(w, metrics)
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8000", nil)
}
