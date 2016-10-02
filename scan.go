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

		fmt.Printf("Detected tuner: %v\n", device.GetTunerType())

		fmt.Printf("Getting gain values...")
		if gains, err := device.GetTunerGains(); err == nil {
			for _, gain := range gains {
				fmt.Printf("* %d\n", gain)
			}
		} else {
			fmt.Println(err.Error())
		}

		fmt.Printf("Setting auto gain mode on ...\n")
		if err := device.SetTunerGainMode(1); err != nil {
			fmt.Println(err.Error())
		}

		device.SetFreqCorrection(52)
		device.SetAgcMode(1)
		device.SetCenterFreq(1090000000)
		device.SetSampleRate(2000000)
		device.ResetBuffer()

		fmt.Printf("Tuner gain: %v\n", device.GetTunerGain())
		fmt.Printf("Device center freq: %v\n", device.GetCenterFreq())
		fmt.Printf("Current sample rate: %v\n", device.GetSampleRate())

		for loop := 0; loop < 1000; loop++ {
			buf, err := device.ReadSync()
			fmt.Printf("Read %d bytes ...\n", len(*buf))
			fmt.Printf("%v", *buf)

			if err != nil {
				panic(err.Error())
			}

		}

		device.Close()
	}

}
