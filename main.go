package main

import (
	"fmt"
	"time"

	"github.com/zingZing1298/go-hardware-monitor/pkgs/hardware"
)

func main() {
	go func() {
		for {
			systemDetails, err := hardware.GetSystemDetails()
			if err != nil {
				fmt.Println(err)
			}

			diskDetails, err := hardware.GetDiskDetails()
			if err != nil {
				fmt.Println(err)
			}

			cpuDetails, err := hardware.GetCPUDetails()
			if err != nil {
				fmt.Println(err)
			}

			netDetails, err := hardware.GetNetworkUsage()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(systemDetails)
			fmt.Println(diskDetails)
			fmt.Println(cpuDetails)
			fmt.Println(netDetails)

			time.Sleep(3 * time.Second)
		}
	}()
	time.Sleep(5 * time.Minute)
}
