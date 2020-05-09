#include "calib3d.h"


void Fisheye_UndistortImage(Mat distorted, Mat undistorted, Mat k, Mat d) {
    cv::fisheye::undistortImage(*distorted, *undistorted, *k, *d);
}

void Fisheye_UndistortImageWithParams(Mat distorted, Mat undistorted, Mat k, Mat d, Mat knew, Size size) {
    cv::Size sz(size.width, size.height);
    cv::fisheye::undistortImage(*distorted, *undistorted, *k, *d, *knew, sz);
}

void InitUndistortRectifyMap(Mat cameraMatrix,Mat distCoeffs,Mat r,Mat newCameraMatrix,Size size,int m1type,Mat map1,Mat map2) {
    cv::Size sz(size.width, size.height);
    cv::initUndistortRectifyMap(*cameraMatrix,*distCoeffs,*r,*newCameraMatrix,sz,m1type,*map1,*map2);
}

Mat GetOptimalNewCameraMatrixWithParams(Mat cameraMatrix,Mat distCoeffs,Size size,double alpha,Size newImgSize,Rect* validPixROI,bool centerPrincipalPoint) {
    cv::Size sz(size.width, size.height);
    cv::Size newSize(newImgSize.width, newImgSize.height);
    cv::Rect rect(validPixROI->x,validPixROI->y,validPixROI->width,validPixROI->height);
    cv::Mat* mat = new cv::Mat(cv::getOptimalNewCameraMatrix(*cameraMatrix,*distCoeffs,sz,alpha,newSize,&rect,centerPrincipalPoint));
    validPixROI->x = rect.x;
    validPixROI->y = rect.y;
    validPixROI->width = rect.width;
    validPixROI->height = rect.height;
    return mat;
}

void Undistort(Mat src, Mat dst, Mat cameraMatrix, Mat distCoeffs, Mat newCameraMatrix) {
    cv::undistort(*src, *dst, *cameraMatrix, *distCoeffs, *newCameraMatrix);
}

bool FindChessboardCorners(Mat image, Size patternSize, Mat corners, int flags) {
    cv::Size sz(patternSize.width, patternSize.height);
    return cv::findChessboardCorners(*image, sz, *corners, flags);
}

void DrawChessboardCorners(Mat image, Size patternSize, Mat corners, bool patternWasFound) {
    cv::Size sz(patternSize.width, patternSize.height);
    cv::drawChessboardCorners(*image, sz, *corners, patternWasFound);
}


// double CalibrateFisheyeCamera(
//     Mat objectPoints,
//     Mat imagePoints,
//     Size imageSize,
//     Mat cameraMatrix,
//     Mat distCoeffs,
//     Mats* outRvecs,
//     Mats* outTvecs,
//     int flags) {

//     std::vector<std::vector<cv::Point3f> > oPts;
//     for (size_t i = 0; i < objectPoints.length; i++) {
//         std::vector<cv::Point3f> op;
//         for(size_t j = 0; j < objectPoints.points[i].length; j++){
//             op.push_back(cv::Point3f(objectPoints.points[i].points[j].x, objectPoints.points[i].points[j].y, objectPoints.points[i].points[j].z));
//         }
//         oPts.push_back(op);
//     }
//     std::vector<std::vector<cv::Point2f> > iPts;
//     for (size_t i = 0; i < imagePoints.length; i++) {
//         std::vector<cv::Point2f> op;
//         for(size_t j = 0; j < imagePoints.points[i].length; j++){
//             op.push_back(cv::Point2f(imagePoints.points[i].points[j].x, imagePoints.points[i].points[j].y));
//         }
//         iPts.push_back(op);
//     }

//     std::vector<cv::Mat> rvecs;
//     std::vector<cv::Mat> tvecs;
//     cv::Size sz(imageSize.width, imageSize.height);
//     double ret = cv::fisheye::calibrate(oPts, iPts, sz, *cameraMatrix, *distCoeffs, rvecs, tvecs, flags);
//     outRvecs->mats = new Mat[rvecs.size()];
//     for (size_t i = 0; i < rvecs.size(); ++i) {
//         outRvecs->mats[i] = new cv::Mat(rvecs[i]);
//     }
//     outRvecs->length = (int)rvecs.size();

//     outTvecs->mats = new Mat[tvecs.size()];
//     for (size_t i = 0; i < tvecs.size(); ++i) {
//         outTvecs->mats[i] = new cv::Mat(tvecs[i]);
//     }
//     outTvecs->length = (int)tvecs.size();

//     return ret;
// }

std::string type2str(int type) {
  std::string r;

  uchar depth = type & CV_MAT_DEPTH_MASK;
  uchar chans = 1 + (type >> CV_CN_SHIFT);

  switch ( depth ) {
    case CV_8U:  r = "8U"; break;
    case CV_8S:  r = "8S"; break;
    case CV_16U: r = "16U"; break;
    case CV_16S: r = "16S"; break;
    case CV_32S: r = "32S"; break;
    case CV_32F: r = "32F"; break;
    case CV_64F: r = "64F"; break;
    default:     r = "User"; break;
  }

  r += "C";
  r += (chans+'0');

  return r;
}

double CalibrateCamera(
    Mats objectPoints,
    Mats imagePoints,
    Size imageSize,
    Mat cameraMatrix,
    Mat distCoeffs,
    Mats* outRvecs,
    Mats* outTvecs,
    int flags) {

    std::ofstream myfile;
    myfile.open ("/tmp/gocv.log");
    myfile << "op len" << ": " << objectPoints.length << ".\n";
    myfile << "ip len" << ": " << imagePoints.length << ".\n";
    myfile.close();

    std::vector<std::vector<cv::Point3f>> oPts;
    for( int m = 0; m < objectPoints.length; ++m ){
        std::vector<cv::Point3f> vec;
        for( int i = 0; i < objectPoints.mats[m]->cols; ++i ){
            for( int j = 0; j < objectPoints.mats[m]->rows; ++j ){
                myfile.open ("/tmp/gocv.log");
                myfile << "ip type" << " type:" << type2str(objectPoints.mats[m]->type()) << ".\n";
                myfile.close();
                vec.push_back(objectPoints.mats[m]->at<cv::Point3f>(j, i, 0));
            }
        }
        oPts.push_back(vec);
    }

    std::vector<std::vector<cv::Point2f> > iPts;
    for( int m = 0; m < imagePoints.length; ++m ){
        std::vector<cv::Point2f> vec;
        for (size_t i = 0; i < imagePoints.mats[m]->cols; i++) {
            for(size_t j = 0; j < imagePoints.mats[m]->rows; j++){
                myfile.open ("/tmp/gocv.log");
                myfile << "ip type i" << i << " j" << j << " : " << type2str(imagePoints.mats[m]->type()) << ".\n";
                myfile.close();
                vec.push_back(imagePoints.mats[m]->at<cv::Point2f>(j, i));
            }
        }
        iPts.push_back(vec);
    }

    cv::Size sz(imageSize.width, imageSize.height);
    std::vector<cv::Mat> rvecs;
    std::vector<cv::Mat> tvecs;

    myfile.open ("/tmp/gocv.log");
    myfile << "opts " << " size:" << oPts.size() << ".\n";
    myfile << "ipts " << " size:" << iPts.size() << ".\n";
    myfile.close();

    double ret = cv::calibrateCamera(oPts, iPts, sz, *cameraMatrix, *distCoeffs, rvecs, tvecs, flags);
    outRvecs->mats = new Mat[rvecs.size()];
    for (size_t i = 0; i < rvecs.size(); ++i) {
        outRvecs->mats[i] = new cv::Mat(rvecs[i]);
    }
    outRvecs->length = (int)rvecs.size();

    outTvecs->mats = new Mat[tvecs.size()];
    for (size_t i = 0; i < tvecs.size(); ++i) {
        outTvecs->mats[i] = new cv::Mat(tvecs[i]);
    }
    outTvecs->length = (int)tvecs.size();

    return ret;
}
