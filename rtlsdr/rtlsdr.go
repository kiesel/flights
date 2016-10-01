package rtlsdr

// #cgo CFLAGS: -I/usr/local/include
// #cgo pkg-config: librtlsdr
// #include "rtl-sdr.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Device struct {
	rtlsdr_dev_t C.struct_rtlsdr_dev_t
}

type Error struct {
	Message string
}

func (this *Error) Error() string {
	return this.Message
}

// RTLSDR_API uint32_t rtlsdr_get_device_count(void);
func Get_device_count() int {
	return int(C.rtlsdr_get_device_count())
}

// RTLSDR_API const char* rtlsdr_get_device_name(uint32_t index);
func Get_device_name(index int) string {
	return C.GoString(C.rtlsdr_get_device_name((C.uint32_t)(index)))
}

// RTLSDR_API int rtlsdr_open(rtlsdr_dev_t **dev, uint32_t index);
func Open(index int) (*Device, error) {
	dev := &Device{}

	result := C.rtlsdr_open(unsafe.Pointer(&dev.rtlsdr_dev_t), (C.uint32_t)(index))
	fmt.Printf("Result = %d", result)

	return dev, nil
}

// RTLSDR_API int rtlsdr_close(rtlsdr_dev_t *dev);
func (this *Device) Close() {
	C.rtlsdr_close(unsafe.Pointer(&this.rtlsdr_dev_t))
}

// RTLSDR_API int rtlsdr_set_tuner_gain(rtlsdr_dev_t *dev, int gain);
func (this *Device) SetTunerGain(gain int) error {
	if 0 != C.rtlsdr_set_tuner_gain(&this.rtlsdr_dev_t, (C.int)(gain)) {
		return Error{"Could not set tuner gain."}
	}

	return nil
}

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
// RTLSDR_API int rtlsdr_get_tuner_gains(rtlsdr_dev_t *dev, int *gains);
// RTLSDR_API int rtlsdr_set_tuner_bandwidth(rtlsdr_dev_t *dev, uint32_t bw);
// RTLSDR_API int rtlsdr_get_tuner_gain(rtlsdr_dev_t *dev);
// RTLSDR_API int rtlsdr_set_tuner_if_gain(rtlsdr_dev_t *dev, int stage, int gain);
// RTLSDR_API int rtlsdr_set_tuner_gain_mode(rtlsdr_dev_t *dev, int manual);
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
