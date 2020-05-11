package main

import (
	"flag"
	"image"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

var (
	chessboardWidth  int
	chessboardHeight int
	imageWidth       int
	imageHeight      int
	squareSize       float64

	srcDir string
	outDir string
)

type CalibrationConfig struct {
	CameraMatrix []byte `yaml:"cameraMatrix"`
	DistCoeffs   []byte `yaml:"distCoeffs"`
}

func main() {
	flag.IntVar(&chessboardWidth, "chessboard-width", 9, "width of chessboard")
	flag.IntVar(&chessboardHeight, "chessboard-height", 6, "height of chessboard")
	flag.IntVar(&imageWidth, "image-width", 1000, "width of input images")
	flag.IntVar(&imageHeight, "image-height", 750, "height of input images")
	flag.Float64Var(&squareSize, "square-size", 7.4, "size of the squares")
	flag.StringVar(&srcDir, "src", "tsst", "directory of input images")
	flag.StringVar(&outDir, "img-out", "out", "directory for output images")
	flag.Parse()

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	var objectPoints []gocv.Mat
	var corners []gocv.Mat

	log.Info("loading images")
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(srcDir, f.Name())
		img := gocv.IMRead(path, gocv.IMReadColor)
		cc := gocv.NewMat()
		if img.Empty() {
			log.Errorf("error: empty image: %s", f.Name())
			continue
		}
		found := gocv.FindChessboardCorners(img, image.Point{chessboardWidth, chessboardHeight}, &cc, gocv.CalibCBAdaptiveThresh|gocv.CalibCBNormalizeImage|gocv.CalibCBFastCheck)
		if !found {
			log.Warnf("chessboard not found: %s", f.Name())
		}
		corners = append(corners, cc)
		imgClone := img.Clone()
		gocv.DrawChessboardCorners(&imgClone, image.Point{chessboardWidth, chessboardHeight}, cc, found)
		imageSize := img.Size()
		printArgs(cc, imageSize)
		ok := gocv.IMWrite(filepath.Join(outDir, f.Name()), imgClone)
		if !ok {
			log.Warnf("could not write output image")
		}
		imgClone.Close()
		src := gocv.NewMatWithSize(1, 6, gocv.MatTypeCV64F)
		setObjectPoints(&src, chessboardWidth, chessboardHeight, float32(squareSize))
		objectPoints = append(objectPoints, src)
	}

	// cleanup
	defer func() {
		for _, f := range objectPoints {
			f.Close()
		}
		for _, c := range corners {
			c.Close()
		}
	}()

	log.Info("calibrating camera")
	cameraMatrix := gocv.NewMat()
	defer cameraMatrix.Close()
	distCoeffs := gocv.NewMat()
	defer distCoeffs.Close()
	flags := 0
	ret := gocv.CalibrateCamera(objectPoints, corners, image.Point{X: imageWidth, Y: imageHeight}, cameraMatrix, distCoeffs, flags)
	printResults(ret, cameraMatrix, distCoeffs)
}

func printArgs(corners gocv.Mat, imageSize []int) {
	for i := 0; i < corners.Rows(); i++ {
		for j := 0; j < corners.Cols(); j++ {
			log.Infof("corner %d %d %f\n", i, j, corners.GetFloatAt(i, j))
		}
	}
	log.Infof("imageSize: %v", imageSize)
}

func printResults(ret float64, cameraMatrix, distCoeffs gocv.Mat) {
	log.Infof("overall RMS re-projection error: %f\n", ret)
	for i := 0; i < cameraMatrix.Rows(); i++ {
		for j := 0; j < cameraMatrix.Cols(); j++ {
			log.Infof("camera Matrix: %d %d %f\n", i, j, cameraMatrix.GetFloatAt(i, j))
		}
	}
	for i := 0; i < distCoeffs.Rows(); i++ {
		for j := 0; j < distCoeffs.Cols(); j++ {
			log.Infof("dist coeff: %d %d %f\n", i, j, distCoeffs.GetFloatAt(i, j))
		}
	}
}

func setObjectPoints(m *gocv.Mat, width, height int, squareSize float32) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			m.SetPoint3fAt(j, i, gocv.Point3f{float32(i) * squareSize, float32(j) * squareSize, 0})
		}
	}
}
