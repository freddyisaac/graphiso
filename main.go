package main

// #cgo CFLAGS: -I/usr/local/include/rav1e
// #cgo LDFLAGS: trampoline.o -L/usr/local/lib -lrav1e
//#include "rav1e.h"
//#include "trampoline.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// could be a generics candidate
func SetConfigValue(raConfig *C.RaConfig, key string, value interface{}) {
	ckey := C.CString("width")
	// will not work for all types - need to restrict to string, bool, int
	cvalue := C.CString(fmt.Sprintf("%v", value))
	C.rav1e_config_parse(raConfig, ckey, cvalue)
	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(cvalue))
}

const (
	width  = 64
	height = 96
	speed  = 9
)

func test_setup() {
	racfg := C.new_rav1e()
	C.t_rav1e_config_default(racfg)

	ret := C.t_rav1e_simple_setup(racfg)
	fmt.Printf("t_rav1e_simple_setup returned %v\n", ret)

	ret = C.t_rav1e_simple_chromaticity(racfg)
	fmt.Printf("t_rav1e_test returned %v\n", ret)

	ret = C.t_rav1e_context_and_frame(racfg)
	fmt.Printf("t_rav1e_context_and_frame returned %v\n", ret)
}

func main() {
	fmt.Printf("hello, world!\n")

	test_setup()

	return

	var raConfig *C.RaConfig

	raConfig = C.rav1e_config_default()

	SetConfigValue(raConfig, "width", width)
	SetConfigValue(raConfig, "height", height)
	SetConfigValue(raConfig, "speed", speed)

	var ret C.int
	ret = C.rav1e_config_set_color_description(raConfig, 2, 2, 2)
	if ret < 0 {
		fmt.Printf("rav1e_config_set_color_description error : %v %T\n", ret, ret)
	}

	var raContext *C.RaContext

	// create a context
	raContext = C.rav1e_context_new(raConfig)

	// create a frame
	var raFrame *C.RaFrame
	fmt.Printf("HELLO1\n")
	raFrame = C.rav1e_frame_new(raContext)
	_ = raFrame
	fmt.Printf("HELLO1\n")

	pixels := make([]uint8, width*height)
	// not the best
	for i := range pixels {
		pixels[i] = 42
	}

	stride := 4
	_ = stride
	l := len(pixels)
	_ = l
	//	C.rav1e_frame_fill_plane(raFrame, 0, (*C.uchar)(&pixels[0]), (C.ulong)(l), (C.long)(stride), 1)

}
