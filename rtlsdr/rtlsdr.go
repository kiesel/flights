package rtlsdr

// #cgo CFLAGS: -I/usr/local/include
// #cgo pkg-config: librtlsdr
// #include "rtl-sdr.h"
// #include <stdlib.h>
//
// rtlsdr_dev_t* rtlsdr_open_helper(uint32_t index) {
//   rtlsdr_dev_t *ptr;
//   int result = rtlsdr_open(&ptr, index);
//
//   if (0 != result) {
//     return NULL;
//   }
//
//   return ptr;
// }
import "C"

import (
	"fmt"
	"unsafe"
)

type Device struct {
	handle *C.struct_rtlsdr_dev
}

type Error struct {
	Message string
}

func (this *Error) Error() string {
	return this.Message
}

type TunerType int

const (
	UNKNOWN = iota
	E4000
	FC0012
	FC0013
	FC2580
	R820T
	R828D
)

func (t TunerType) String() string {
	switch {
	case t == E4000:
		return "E4000"
	case t == FC0012:
		return "FC0012"
	case t == FC0013:
		return "FC0013"
	case t == FC2580:
		return "FC2580"
	case t == R820T:
		return "R820T"
	case t == R828D:
		return "R828D"
	}

	return "UNKNOWN"
}

// RTLSDR_API uint32_t rtlsdr_get_device_count(void);
func GetDeviceCount() int {
	return int(C.rtlsdr_get_device_count())
}

// RTLSDR_API const char* rtlsdr_get_device_name(uint32_t index);
func GetDeviceName(index int) string {
	return C.GoString(C.rtlsdr_get_device_name((C.uint32_t)(index)))
}

// RTLSDR_API int rtlsdr_open(rtlsdr_dev_t **dev, uint32_t index);
func Open(index int) (*Device, error) {
	dev := C.rtlsdr_open_helper(((C.uint32_t)(index)))

	if unsafe.Pointer(dev) == unsafe.Pointer(nil) {
		return nil, &Error{"Unable to open rtlsdr device."}
	}

	return &Device{handle: dev}, nil
}

// RTLSDR_API int rtlsdr_close(rtlsdr_dev_t *dev);
func (this *Device) Close() {
	C.rtlsdr_close((*C.struct_rtlsdr_dev_t)(unsafe.Pointer(this.handle)))
}

// RTLSDR_API int rtlsdr_set_tuner_gain(rtlsdr_dev_t *dev, int gain);
func (this *Device) SetTunerGain(gain int) error {
	if 0 != C.rtlsdr_set_tuner_gain((*C.struct_rtlsdr_dev_t)(unsafe.Pointer(this.handle)), (C.int)(gain)) {
		return &Error{"Could not set tuner gain."}
	}

	return nil
}

// RTLSDR_API int rtlsdr_set_tuner_gain_mode(rtlsdr_dev_t *dev, int manual);
func (this *Device) SetTunerGainMode(manual int) error {
	if res := C.rtlsdr_set_tuner_gain_mode((*C.struct_rtlsdr_dev_t)(unsafe.Pointer(this.handle)), ((C.int)(manual))); res != 0 {
		return &Error{"Unable to set tuner gain mode."}
	}

	return nil
}

// RTLSDR_API int rtlsdr_get_tuner_gains(rtlsdr_dev_t *dev, int *gains);
func (this *Device) GetTunerGains() ([]int, error) {
	// First, ask how many gain values will be returned
	numberOfGains := (int)(C.rtlsdr_get_tuner_gains(this.handle, (*C.int)(unsafe.Pointer(nil))))
	fmt.Printf("Gain numer = %d\n", numberOfGains)

	gains := make([]int, numberOfGains)
	result := (C.rtlsdr_get_tuner_gains(this.handle, (*C.int)(unsafe.Pointer(&gains[0]))))
	fmt.Printf("Gain result = %d\n", result)

	if result <= 0 {
		return []int{}, &Error{"Could not retrieve gain values."}
	}

	return gains, nil
}

// RTLSDR_API int rtlsdr_get_tuner_gain(rtlsdr_dev_t *dev);
func (this *Device) GetTunerGain() int {
	gain := C.rtlsdr_get_tuner_gain(this.handle)
	// if gain == 0 {
	//   return 0, &Error{"Could not get tuner gain."}
	// }

	// return (int)(gain), nil
	return (int)(gain)
}

// RTLSDR_API int rtlsdr_set_tuner_if_gain(rtlsdr_dev_t *dev, int stage, int gain);
func (this *Device) SetTunerIfGain(stage, gain int) error {
	if result := C.rtlsdr_set_tuner_if_gain(this.handle, (C.int)(stage), (C.int)(gain)); result != 0 {
		return &Error{"Could not set intermediate freq. gain"}
	}

	return nil
}

// RTLSDR_API int rtlsdr_set_freq_correction(rtlsdr_dev_t *dev, int ppm);
func (this *Device) SetFreqCorrection(ppm int) error {
	result := C.rtlsdr_set_freq_correction(this.handle, (C.int)(ppm))
	if result != 0 {
		return &Error{"Could not set freq correction"}
	}

	return nil
}

// RTLSDR_API int rtlsdr_get_freq_correction(rtlsdr_dev_t *dev);
func (this *Device) GetFreqCorrection() int {
	return (int)(C.rtlsdr_get_freq_correction(this.handle))
}

// RTLSDR_API enum rtlsdr_tuner rtlsdr_get_tuner_type(rtlsdr_dev_t *dev);
func (this *Device) GetTunerType() TunerType {
	return (TunerType)(C.rtlsdr_get_tuner_type(this.handle))
}

// RTLSDR_API int rtlsdr_set_center_freq(rtlsdr_dev_t *dev, uint32_t freq);
func (this *Device) SetCenterFreq(freq uint32) error {
	result := C.rtlsdr_set_center_freq(this.handle, (C.uint32_t)(freq))
	if result != 0 {
		return &Error{"Cannot set center freq"}
	}

	return nil
}

// RTLSDR_API uint32_t rtlsdr_get_center_freq(rtlsdr_dev_t *dev);
func (this *Device) GetCenterFreq() uint32 {
	return (uint32)(C.rtlsdr_get_center_freq(this.handle))
}

// RTLSDR_API int rtlsdr_set_agc_mode(rtlsdr_dev_t *dev, int on);
func (this *Device) SetAgcMode(on int) error {
	result := C.rtlsdr_set_agc_mode(this.handle, (C.int)(on))

	if 0 != result {
		return &Error{"Could not set agc mode."}
	}

	return nil
}

// RTLSDR_API int rtlsdr_set_sample_rate(rtlsdr_dev_t *dev, uint32_t rate);
func (this *Device) SetSampleRate(rate uint32) error {
	result := C.rtlsdr_set_sample_rate(this.handle, (C.uint32_t)(rate))

	if 0 != result {
		return &Error{"Could not set sample rate"}
	}

	return nil
}

// RTLSDR_API uint32_t rtlsdr_get_sample_rate(rtlsdr_dev_t *dev);
func (this *Device) GetSampleRate() uint32 {
	return (uint32)(C.rtlsdr_get_sample_rate(this.handle))
}

// RTLSDR_API int rtlsdr_reset_buffer(rtlsdr_dev_t *dev);
func (this *Device) ResetBuffer() {
	C.rtlsdr_reset_buffer(this.handle)
}

// RTLSDR_API int rtlsdr_read_sync(rtlsdr_dev_t *dev, void *buf, int len, int *n_read);
func (this *Device) ReadSync(buflen int) (*[]byte, error) {
	fmt.Printf("reading up to %d bytes\n", buflen)
	buf := make([]byte, buflen)
	n_read := (C.int)(0)
	result := C.rtlsdr_read_sync(this.handle, unsafe.Pointer(&buf[0]), (C.int)(buflen), &n_read)

	fmt.Printf("read %d bytes\n", (int)(n_read))
	if result != 0 {
		return nil, &Error{"Could not read bytes from device"}
	}

	return &buf, nil
}

// RTLSDR_API int rtlsdr_get_device_usb_strings(uint32_t index,
// RTLSDR_API int rtlsdr_get_index_by_serial(const char *serial);
// RTLSDR_API int rtlsdr_set_xtal_freq(rtlsdr_dev_t *dev, uint32_t rtl_freq,
// RTLSDR_API int rtlsdr_get_xtal_freq(rtlsdr_dev_t *dev, uint32_t *rtl_freq,
// RTLSDR_API int rtlsdr_get_usb_strings(rtlsdr_dev_t *dev, char *manufact,
// RTLSDR_API int rtlsdr_write_eeprom(rtlsdr_dev_t *dev, uint8_t *data,
// RTLSDR_API int rtlsdr_read_eeprom(rtlsdr_dev_t *dev, uint8_t *data,
// RTLSDR_API int rtlsdr_set_tuner_bandwidth(rtlsdr_dev_t *dev, uint32_t bw);
// RTLSDR_API int rtlsdr_set_testmode(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_set_direct_sampling(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_get_direct_sampling(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_offset_tuning(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_get_offset_tuning(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_wait_async(rtlsdr_dev_t *dev, rtlsdr_read_async_cb_t cb, void *ctx);
// RTLSDR_API int rtlsdr_read_async(rtlsdr_dev_t *dev,
// RTLSDR_API int rtlsdr_cancel_async(rtlsdr_dev_t *dev);
