package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zingZing1298/go-hardware-monitor/pkgs/hardware"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
}

func MakeNewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
	}
	// Setting up default router to get server static files in given location.
	s.mux.Handle("/", http.FileServer(http.Dir("./templates")))
	return s
}

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
	// time.Sleep(5 * time.Minute)\
	// Run server
	srv := MakeNewServer()
	err := http.ListenAndServe(":8000", &srv.mux)

	if err != nil {
		fmt.Println("Server crash...\n", err)
		os.Exit(1)
	}

}
