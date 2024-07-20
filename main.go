package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type Metrics struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	NetworkSent uint64
	NetworkRecv	uint64
}

func getMetrics() Metrics {
	// Get CPU Usage
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory Usage
	vmStat, _ := mem.VirtualMemory()

	// Get disk usage
	diskUsage, _ := disk.Usage("/")

	//get network Usage
	netStats,_ := net.IOCounters(false)
	return Metrics{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: vmStat.UsedPercent,
		DiskUsage:   diskUsage.UsedPercent,
		NetworkSent: netStats[0].BytesSent,
		NetworkRecv: netStats[0].BytesRecv,

	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := getMetrics()
	tmpl := template.Must(template.ParseFiles("templates/metrics.html"))
	tmpl.Execute(w, metrics)
}

func redirectHandler(w http.ResponseWriter, r *http.Request){
	http.Redirect(w,r, "/metrics",http.StatusMovedPermanently)

}


func main() {
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/metrics", metricsHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil{
		log.Fatalf("Could not start server:  %s \n", err)
	}
}
