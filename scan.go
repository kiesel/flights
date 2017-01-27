package main

import (
	"fmt"
	"github.com/kiesel/flights/rtlsdr"
  "github.com/inpursuit/manchester"
)

const (
  PREAMBLE_LEN = 16

  long_frame = 112
  short_frame = 56
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

		var channel = make(chan *[]byte)

		go processReceivedData(channel)

		for loop := 0; loop < 50; loop++ {
			buf, err := device.ReadSync(16 * 16384)

			if err != nil {
				panic(err.Error())
			}

			channel <- buf
		}

		device.Close()
	}
}

func processReceivedData(queue chan *[]byte) {
	for {
		buf := <-queue
		fmt.Printf("Received %d bytes for processing\n", len(*buf))

    // Convert received data into uint16
    mag := make([]uint16, len(*buf) / 2)
    for i, j := 0, 0; i < len(*buf); i, j = i+2, j+1 {
      mag[j]= (uint16)((*buf)[i] * (*buf)[i] + (*buf)[i+1] * (*buf)[i+1])
    }

    // Apply manchester
    manchester.Manchester(mag)
    
    // Extract messages
    extractMessages(mag)
	}
}

func extractMessages(buf []uint16) {
  var frame_len, data_i, index int
  var shift uint8
  var adsb_frame = make([]int, 14)

  for i := 0; i 

  // fmt.Printf("%v\n", buf)

  for i, numBytes := 0, len(buf); i < numBytes; i++ {
    if buf[i] > 1 {
      continue
    }

    frame_len = long_frame
    data_i = 0

    for ; i < frame_len && buf[i] <= 1 && data_i < frame_len; i, data_i = i + 1, data_i + 1 {
      if buf[i] != 0 {
        index = data_i / 8
        shift = (uint8)(7 - (data_i % 8))
        adsb_frame[index] |= (1 << shift)
        // fmt.Printf("i = %d, data_i = %d data = %d, index = %d, adsb_frame = %d\n", i, data_i, buf[i], index, adsb_frame[index])
      }

      if data_i == 7 {
        if adsb_frame[0] == 0 {
          break;
        }

        if (adsb_frame[0] & 0x80) != 0 {
          frame_len = long_frame
        } else {
          frame_len = short_frame;
        }
      }
    }

    if data_i < (frame_len - 1) {
      // fmt.Printf("%d / %d\n", data_i, frame_len - 1)
      continue;
    }

    displayMessage(adsb_frame)
  }
}

func displayMessage(frame []int) {
  df := (frame[0] >> 3) & 0x1f

  for i := 0; i < ((len(frame) + 7) / 8); i++ {
    fmt.Printf(" 0x%h", frame[i])
  }
  fmt.Println()

  fmt.Printf("DF=%d CA=%d\n", df, frame[0] & 0x07)
  fmt.Printf("ICAO=%06x\n", frame[1] << 16 | frame[2] << 8 | frame[3])
  if len(frame) <= short_frame {
    return
  }
}
