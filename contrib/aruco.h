#ifndef _OPENCV3_ARUCO_H_
#define _OPENCV3_ARUCO_H_

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
#include <opencv2/aruco.hpp>
extern "C" {
#endif

#include "../core.h"

#ifdef __cplusplus
typedef cv::Ptr<cv::aruco::Dictionary>* ArucoDictionary;
#else
typedef void* ArucoDictionary;
#endif

typedef struct ArucoCorner {
    int id;
    Point2f p1;
    Point2f p2;
    Point2f p3;
    Point2f p4;
} ArucoCorner;

typedef struct ArucoCorners {
    ArucoCorner* corners;
    int length;
} ArucoCorners;

typedef struct ArucoDetectedMarkers {
    ArucoCorners* markers;
    ArucoCorners* rejected;
} ArucoDetectedMarkers;

typedef struct ArucoPoseVecs {
    Mats* rvecs;
    Mats* tvecs;
} ArucoPoseVecs;

ArucoDictionary ArucoDictionary_Create(Mat bytesList, int markerSize, int maxcorr);
ArucoDictionary ArucoDictionaryPredefined(int t);
void ArucoDictionary_DrawMarker(ArucoDictionary d, int id, int sidePixels, Mat out);
void ArucoDictionary_Close(ArucoDictionary b);
ArucoDetectedMarkers ArucoDetectMarkers(Mat img, ArucoDictionary d);
void ArucoDrawAxis(Mat dst, Mat cameraMatrix, Mat distCoeffs, Mat rvec, Mat tvec, float length);
void ArucoEstimatePoseSingleMarkers(ArucoCorners* corners, float markerLength, Mat cameraMatrix, Mat distCoeffs, struct Mats* rvc, struct Mats* tvec);
void ArucoCorners_Close(ArucoCorners* c);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV3_ARUCO_H_
