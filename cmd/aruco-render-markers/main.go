package main

import (
	"fmt"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func main() {

	// create dict & draw marker
	d2 := contrib.NewArucoPredefinedDictionary(contrib.ArucoPredefinedDict6X6_250)
	defer d2.Close()

	// todo: get number of markers
	// print markers
	for i := 0; i < 12; i++ {
		dst := gocv.NewMat()
		d2.DrawMarker(i, 500, &dst)
		gocv.IMWrite(fmt.Sprintf("./aruco-%d.png", i), dst)
		dst.Close()
	}
}
