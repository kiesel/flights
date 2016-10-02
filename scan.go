package main

import (
	"fmt"
	"github.com/kiesel/flights/rtlsdr"
)

func main() {
	var device *rtlsdr.Device
	var err error

	for i := 0; i < rtlsdr.GetDeviceCount(); i++ {
		fmt.Printf("Found device '%s'\n", rtlsdr.GetDeviceName(i))
		if device, err = rtlsdr.Open(i); err != nil {
			panic(err)
		}

		fmt.Printf("Getting gain values...")
		if gains, err := device.GetTunerGains(); err != nil {
			fmt.Printf("%v", gains)
			for _, gain := range gains {
				fmt.Printf("* %d\n", gain)
			}
		} else {
			fmt.Printf(err.Error())
		}

		device.Close()
	}

}
