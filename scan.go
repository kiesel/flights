package main

import (
	"fmt"
	"github.com/kiesel/flights/rtlsdr"
)

func main() {
	var device *rtlsdr.Device
	var err rtlsdr.Error

	for i := 0; i < rtlsdr.Get_device_count(); i++ {
		fmt.Printf("Found device '%s'\n", rtlsdr.Get_device_name(i))
		if device, err = rtlsdr.Open(i); err != nil {
			panic(err)
		}

		device.Close()
	}

}
