package darknet

// #include <stdlib.h>
//
// #cgo LDFLAGS: -ldarknet
// void create_detector(int xpu, char *cfgfile, char *weightfile, void **handle);
// void forward_detector(void *handle, unsigned char *CHW, int c, int h, int w, float thresh, float hier_thresh, float nms, float **out, unsigned short *out_len);
// void free_detector(void *handle);
import "C"
import (
	"unsafe"
)

// Detector represents a darknet detector
type Detector struct {
	handle unsafe.Pointer
}

// Detection includes the index of the label, its probability, and the position of the left-top and right-bottom corner of the bounding box, in the input coordiante, i.e. [0, W) x [0, H).
type Detection struct {
	LabelIdx    uint
	Probability float32
	Left        uint
	Top         uint
	Right       uint
	Bottom      uint
}

// NewDetector creates a new Detector. `xpu` is the index of GPU, or -1 indicating CPU.
func NewDetector(xpu int, symbolPath string, paramPath string) (detector *Detector, err error) {
	detector = &Detector{}

	// do not check the returned error.
	// refer to https://github.com/hxhxhx88/go-mxnet-predictor/blob/master/mxnet/predictor.go#L86
	C.create_detector(
		C.int(xpu),
		C.CString(symbolPath),
		C.CString(paramPath),
		&detector.handle,
	)
	return
}

// Detect an image. The image should be given in CHW order.
func (d *Detector) Detect(ps []uint8, c int, h int, w int, thres float32, hierThres float32, nmsThres float32) (detections []Detection, err error) {
	var cDets *C.float
	var cLen C.ushort

	// do not check the returned error.
	// refer to https://github.com/hxhxhx88/go-mxnet-predictor/blob/master/mxnet/predictor.go#L86
	C.forward_detector(
		d.handle,
		(*C.uchar)(unsafe.Pointer(&ps[0])),
		C.int(c),
		C.int(h),
		C.int(w),
		C.float(thres),
		C.float(hierThres),
		C.float(nmsThres),
		&cDets,
		&cLen,
	)

	// c array to go
	dets := (*[1 << 16]float32)(unsafe.Pointer(cDets))[:cLen:cLen]

	// `dets` consists of several [label, prob, left, top, right, bottom]s
	for i := 0; i < len(dets); i += 6 {
		var det Detection
		det.LabelIdx = uint(dets[i])
		det.Probability = dets[i+1]
		det.Left = uint(dets[i+2])
		det.Top = uint(dets[i+3])
		det.Right = uint(dets[i+4])
		det.Bottom = uint(dets[i+5])
		detections = append(detections, det)
	}

	// free mem we created before return, go gc won't do that for us
	C.free(unsafe.Pointer(cDets))

	return
}

// Free ...
func (d *Detector) Free() (err error) {
	// do not check the returned error.
	// refer to https://github.com/hxhxhx88/go-mxnet-predictor/blob/master/mxnet/predictor.go#L86
	C.free_detector(d.handle)
	return
}
