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

		// var channel = make(chan *[]byte)
		// go processReceivedData(channel)

		for loop := 0; loop < 10; loop++ {
			fmt.Printf("Now reading ...\n")
			buf, err := device.ReadSync(16)

			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("Handing off data of len %d\n", len(*buf))
			// channel <- buf
		}

		device.Close()
	}
}

func processReceivedData(queue chan *[]byte) {
	for {
		fmt.Println("Waiting for data ...")
		buf := <-queue
		fmt.Printf("Received %d bytes for processing\n", len(*buf))
	}
}
