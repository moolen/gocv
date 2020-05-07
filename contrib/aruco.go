package contrib

/*
#include <stdlib.h>
#include "aruco.h"
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

// ArucoDictionary is a set of markers. It contains the inner codification
type ArucoDictionary struct {
	// C.ArucoDictionary
	p unsafe.Pointer
}

// ArucoPredefinedDict are predefined markers dictionaries/sets
// Each dictionary indicates the number of bits and the number of markers contained.
type ArucoPredefinedDict int

const (
	ArucoPredefinedDict4X4_50 ArucoPredefinedDict = iota
	ArucoPredefinedDict4X4_100
	ArucoPredefinedDict4X4_250
	ArucoPredefinedDict4X4_1000
	ArucoPredefinedDict5X5_50
	ArucoPredefinedDict5X5_100
	ArucoPredefinedDict5X5_250
	ArucoPredefinedDict5X5_1000
	ArucoPredefinedDict6X6_50
	ArucoPredefinedDict6X6_100
	ArucoPredefinedDict6X6_250
	ArucoPredefinedDict6X6_1000
	ArucoPredefinedDict7X7_50
	ArucoPredefinedDict7X7_100
	ArucoPredefinedDict7X7_250
	ArucoPredefinedDict7X7_1000
	// ArucoPredefinedDictArucoOriginal are the standard ArUco Library Markers.
	// 1024 markers, 5x5 bits, 0 minimum distance
	ArucoPredefinedDictArucoOriginal
	// ArucoPredefinedDictAprilTag16h5 is 4x4 bits, minimum hamming distance between any two codes = 5, 30 codes
	ArucoPredefinedDictAprilTag16h5
	// ArucoPredefinedDictAprilTag25h9 is 5x5 bits, minimum hamming distance between any two codes = 9, 35 codes
	ArucoPredefinedDictAprilTag25h9
	// ArucoPredefinedDictAprilTag36h10 is 6x6 bits, minimum hamming distance between any two codes = 10, 2320 codes
	ArucoPredefinedDictAprilTag36h10
	// ArucoPredefinedDictAprilTag36h11 is 6x6 bits, minimum hamming distance between any two codes = 11, 587 codes
	ArucoPredefinedDictAprilTag36h11
)

// NewArucoDictionary returns a empty Dictionary
func NewArucoDictionary(byteList *gocv.Mat, markerSize, maxCorr int) *ArucoDictionary {
	return &ArucoDictionary{p: unsafe.Pointer(C.ArucoDictionary_Create(C.Mat(byteList.Ptr()), C.int(markerSize), C.int(maxCorr)))}
}

// NewArucoPredefinedDictionary returns a new dictionary from a preset
func NewArucoPredefinedDictionary(preset ArucoPredefinedDict) *ArucoDictionary {
	return &ArucoDictionary{p: unsafe.Pointer(C.ArucoDictionaryPredefined(C.int(preset)))}
}

// DrawMarker renders marker with specified id into the output image
// this fails if id does not exist in dict or sidePixels < (markerSize + 2*borderBits)
func (d *ArucoDictionary) DrawMarker(id int, sidePixels int, dst *gocv.Mat) {
	C.ArucoDictionary_DrawMarker((C.ArucoDictionary)(d.p), C.int(id), C.int(sidePixels), (C.Mat)(dst.Ptr()))
}

// Close frees the allocated resources
func (d *ArucoDictionary) Close() error {
	C.ArucoDictionary_Close((C.ArucoDictionary)(d.p))
	d.p = nil
	return nil
}

// ArucoCorner is a marker with an id and four points
// that describe the position of the marker within the image
type ArucoCorner struct {
	ID int
	P1 gocv.Point2f
	P2 gocv.Point2f
	P3 gocv.Point2f
	P4 gocv.Point2f
}

type ArucoCorners struct {
	p       unsafe.Pointer
	corners []ArucoCorner
}

// TODO: implement detectorParameters
// https://docs.opencv.org/trunk/d1/dcd/structcv_1_1aruco_1_1DetectorParameters.html

// ArucoDetectMarkers detects markers in img which are defined in dict
// it returns the markers found and rejected candidates with their respective coordinates
// Markers and rejected candidates must be freed after use.
func ArucoDetectMarkers(img *gocv.Mat, dict *ArucoDictionary) (*ArucoCorners, *ArucoCorners) {
	logrus.Printf("args: %v | %v", (C.Mat)(img.Ptr()), (C.ArucoDictionary)(dict.p))
	ret := C.ArucoDetectMarkers((C.Mat)(img.Ptr()), (C.ArucoDictionary)(dict.p))
	return getDetectedMarkers(ret)
}

func getDetectedMarkers(ret C.ArucoDetectedMarkers) (*ArucoCorners, *ArucoCorners) {
	cFoundMarker := ret.markers.corners
	markerLen := int(ret.markers.length)
	cRejectedMarker := ret.rejected.corners
	rejectedLength := int(ret.rejected.length)
	foundMarker := getArucoCorners(unsafe.Pointer(ret.markers), unsafe.Pointer(cFoundMarker), markerLen)
	rejectedMarker := getArucoCorners(unsafe.Pointer(ret.rejected), unsafe.Pointer(cRejectedMarker), rejectedLength)
	return foundMarker, rejectedMarker
}

func getArucoCorners(arrPtr unsafe.Pointer, data unsafe.Pointer, len int) *ArucoCorners {
	hdr := reflect.SliceHeader{
		Data: uintptr(data),
		Len:  len,
		Cap:  len,
	}
	s := *(*[]C.ArucoCorner)(unsafe.Pointer(&hdr))
	keys := make([]ArucoCorner, len)
	for i, r := range s {
		keys[i] = ArucoCorner{
			ID: int(r.id),
			P1: gocv.Point2f{float32(r.p1.x), float32(r.p1.y)},
			P2: gocv.Point2f{float32(r.p2.x), float32(r.p2.y)},
			P3: gocv.Point2f{float32(r.p3.x), float32(r.p3.y)},
			P4: gocv.Point2f{float32(r.p4.x), float32(r.p4.y)},
		}
	}
	return &ArucoCorners{
		p:       arrPtr,
		corners: keys,
	}
}

func (c *ArucoCorners) Close() {
	C.ArucoCorners_Close((*C.ArucoCorners)(c.p))
}

func (c *ArucoCorners) Range() <-chan ArucoCorner {
	ch := make(chan ArucoCorner)
	go func() {
		for _, corner := range c.corners {
			ch <- corner
		}
		close(ch)
	}()
	return ch
}

func (c *ArucoCorners) Len() int {
	return len(c.corners)
}

// ArucoDrawAxis draws the axis on the dst image using the camera calibration matrices
func ArucoDrawAxis(dst *gocv.Mat, cameraMatrix gocv.Mat, distCoeffs gocv.Mat, rvec gocv.Mat, tvec gocv.Mat, length float32) {
	C.ArucoDrawAxis(C.Mat(dst.Ptr()), C.Mat(cameraMatrix.Ptr()), C.Mat(distCoeffs.Ptr()), C.Mat(rvec.Ptr()), C.Mat(tvec.Ptr()), C.float(length))
}

// ArucoEstimatePoseSingleMarkers does pose estimation for single markers
func ArucoEstimatePoseSingleMarkers(c *ArucoCorners, markerLength float32, cameraMatrix gocv.Mat, distCoeffs gocv.Mat) ([]gocv.Mat, []gocv.Mat) {
	cRVecs := C.struct_Mats{}
	cTVecs := C.struct_Mats{}
	C.ArucoEstimatePoseSingleMarkers(
		(*C.ArucoCorners)(unsafe.Pointer(c.p)),
		C.float(markerLength),
		C.Mat(cameraMatrix.Ptr()),
		C.Mat(distCoeffs.Ptr()),
		&(cRVecs),
		&(cTVecs))

	return gocv.NewMatsFromPtr(unsafe.Pointer(&cRVecs)),
		gocv.NewMatsFromPtr(unsafe.Pointer(&cTVecs))
}

// https://docs.opencv.org/master/d9/d6a/group__aruco.html
// aruco::estimatePoseBoard()
