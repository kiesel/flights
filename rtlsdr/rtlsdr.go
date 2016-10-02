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
// RTLSDR_API int rtlsdr_set_tuner_if_gain(rtlsdr_dev_t *dev, int stage, int gain);

// RTLSDR_API int rtlsdr_get_device_usb_strings(uint32_t index,
// RTLSDR_API int rtlsdr_get_index_by_serial(const char *serial);
// RTLSDR_API int rtlsdr_set_xtal_freq(rtlsdr_dev_t *dev, uint32_t rtl_freq,
// RTLSDR_API int rtlsdr_get_xtal_freq(rtlsdr_dev_t *dev, uint32_t *rtl_freq,
// RTLSDR_API int rtlsdr_get_usb_strings(rtlsdr_dev_t *dev, char *manufact,
// RTLSDR_API int rtlsdr_write_eeprom(rtlsdr_dev_t *dev, uint8_t *data,
// RTLSDR_API int rtlsdr_read_eeprom(rtlsdr_dev_t *dev, uint8_t *data,
// RTLSDR_API int rtlsdr_set_center_freq(rtlsdr_dev_t *dev, uint32_t freq);
// RTLSDR_API uint32_t rtlsdr_get_center_freq(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_freq_correction(rtlsdr_dev_t *dev, int ppm);
// RTLSDR_API int rtlsdr_get_freq_correction(rtlsdr_dev_t *dev);
// RTLSDR_API enum rtlsdr_tuner rtlsdr_get_tuner_type(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_tuner_bandwidth(rtlsdr_dev_t *dev, uint32_t bw);
// RTLSDR_API int rtlsdr_set_sample_rate(rtlsdr_dev_t *dev, uint32_t rate);
// RTLSDR_API uint32_t rtlsdr_get_sample_rate(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_testmode(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_set_agc_mode(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_set_direct_sampling(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_get_direct_sampling(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_offset_tuning(rtlsdr_dev_t *dev, int on);
// RTLSDR_API int rtlsdr_get_offset_tuning(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_reset_buffer(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_read_sync(rtlsdr_dev_t *dev, void *buf, int len, int *n_read);
// RTLSDR_API int rtlsdr_wait_async(rtlsdr_dev_t *dev, rtlsdr_read_async_cb_t cb, void *ctx);
// RTLSDR_API int rtlsdr_read_async(rtlsdr_dev_t *dev,
// RTLSDR_API int rtlsdr_cancel_async(rtlsdr_dev_t *dev);
