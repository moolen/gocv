package main

import (
	"image"
	"image/color"

	log "github.com/sirupsen/logrus"
	gocv "gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func main() {
	img := gocv.IMRead("./aruco_example.jpg", gocv.IMReadColor)
	defer img.Close()
	if img.Empty() {
		log.Fatalf("empty image: loading probably failed.")
	}
	d2 := contrib.NewArucoPredefinedDictionary(contrib.ArucoPredefinedDict6X6_250)
	defer d2.Close()

	// detect markers
	corners, rejected := contrib.ArucoDetectMarkers(&img, d2)
	log.Infof("corners: %d marker %v\n", corners.Len(), corners)
	defer corners.Close()
	defer rejected.Close()

	// draw markers
	for m := range corners.Range() {
		drawRect(&img, m, color.RGBA{0, 255, 0, 1})
	}
	//for m := range rejected.Range() {
	//drawRect(&img, m, color.RGBA{255, 0, 0, 1})
	//}
	cameraMatrix := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV64F)
	defer cameraMatrix.Close()
	distCoeffs := gocv.NewMatWithSize(5, 1, gocv.MatTypeCV64F)
	defer distCoeffs.Close()
	rvec := gocv.NewMatWithSize(12, 12, gocv.MatTypeCV64F)
	defer rvec.Close()
	tvec := gocv.NewMatWithSize(12, 12, gocv.MatTypeCV64F)
	defer tvec.Close()
	rvecs, tvecs := contrib.ArucoEstimatePoseSingleMarkers(corners, 3.3, cameraMatrix, distCoeffs)

	log.Infof("%v, %v", rvecs, tvecs)
	for i := 0; i < corners.Len(); i++ {
		contrib.ArucoDrawAxis(&img, cameraMatrix, distCoeffs, rvecs[i], tvecs[i], 124.4)
	}

	gocv.IMWrite("./aruco_example_detected.jpg", img)
}

func drawRect(img *gocv.Mat, m contrib.ArucoCorner, c color.RGBA) {
	gocv.Line(img,
		image.Point{int(m.P1.X), int(m.P1.Y)}, image.Point{int(m.P2.X), int(m.P2.Y)},
		c, 3,
	)
	gocv.Line(img,
		image.Point{int(m.P2.X), int(m.P2.Y)}, image.Point{int(m.P1.X), int(m.P1.Y)},
		c, 3,
	)
	gocv.Line(img,
		image.Point{int(m.P3.X), int(m.P3.Y)}, image.Point{int(m.P2.X), int(m.P2.Y)},
		c, 3,
	)
	gocv.Line(img,
		image.Point{int(m.P1.X), int(m.P1.Y)}, image.Point{int(m.P4.X), int(m.P4.Y)},
		c, 3,
	)
}
