package contrib

import (
	"image"
	"image/color"
	"testing"

	"gocv.io/x/gocv"
)

func TestAruco(t *testing.T) {
	dst := gocv.NewMat()
	defer dst.Close()

	// create dict & draw marker
	d2 := NewArucoPredefinedDictionary(ArucoPredefinedDict6X6_250)
	defer d2.Close()
	d2.DrawMarker(12, 400, &dst)
	gocv.IMWrite("../images/aruco.png", dst)

	// detect markers
	img := gocv.IMRead("../images/aruco_example.jpg", gocv.IMReadColor)
	found, rejected := ArucoDetectMarkers(&img, d2)
	expectedMarker := 6
	if len(found) != expectedMarker {
		t.Errorf("invalid number of markers found. expected %d, found: %d - %v", expectedMarker, len(found), found)
	}

	// draw markers (debugging)
	for _, m := range found {
		drawRect(&img, m, color.RGBA{0, 255, 0, 1})
	}
	for _, m := range rejected {
		drawRect(&img, m, color.RGBA{255, 0, 0, 1})
	}
	gocv.IMWrite("../images/aruco_example_detect.jpg", img)
}

func drawRect(img *gocv.Mat, m ArucoDetectedMarker, c color.RGBA) {
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
