package pkgs

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func GetCPUDetails() (string, error) {
	cpuStat, err := cpu.Info()

	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("CPU: %s \nCores: %d", cpuStat[0].ModelName, len(cpuStat))
	return output, nil
}

func GetDiskDetails() (string, error) {
	diskStat, err := disk.Usage("/")

	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("Total Disk Space: %d \nFree Space: %d", diskStat.Total, diskStat.Free)
	return output, nil
}

func GetSystemDetails() (string, error) {
	runTimeOS := runtime.GOOS
	vimStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("HostName: %s\nTotal Memory: %d\nUsed Memory: %d\nOS: %s", hostStat.Hostname, vimStat.Total, vimStat.Used, runTimeOS)
	return output, nil
}

func GetNetworkUsage() (string, error) {
	netStat, err := net.IOCounters(false)
	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("Bytes Sent: %d\nBytes Received: %d", netStat[0].BytesSent, netStat[0].BytesRecv)
	return output, nil
}
