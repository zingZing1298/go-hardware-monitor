package main

import (
	"fmt"
	"pkgs/hardware"
	
)

func main(){
	go func()  {
		for{
			systemDetails, err:= hardware.GetSystemDetails()
			if err!=nil{
				fmt.Println(err)
			}
			






		}
	}
}