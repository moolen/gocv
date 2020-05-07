#include "aruco.h"

ArucoDictionary ArucoDictionary_Create(Mat bytesList, int markerSize, int maxcorr) {
    return new cv::Ptr<cv::aruco::Dictionary>(new cv::aruco::Dictionary(*bytesList, markerSize, maxcorr));
}

ArucoDictionary ArucoDictionaryPredefined(int t) {
    return new cv::Ptr<cv::aruco::Dictionary>(cv::aruco::getPredefinedDictionary(t));
}

void ArucoDictionary_Close(ArucoDictionary d) {
    delete d;
}

void ArucoDictionary_DrawMarker(ArucoDictionary d, int id, int sidePixels, Mat dst){
    (*d)->drawMarker(id, sidePixels, *dst);
}

ArucoDetectedMarkers ArucoDetectMarkers(Mat img, ArucoDictionary d){
    cv::Ptr<cv::aruco::DetectorParameters> dp(new cv::aruco::DetectorParameters());
    std::vector<int> markerIds;
    std::vector<std::vector<cv::Point2f>> markerCorners;
    std::vector<std::vector<cv::Point2f>> markerRejected;

    // todo: parameterize DetectorParameter
    dp->minDistanceToBorder = 0;
    dp->adaptiveThreshWinSizeMax = 400;

    cv::aruco::detectMarkers(*img, *d, markerCorners,
        markerIds, dp, markerRejected);

    ArucoCorners* markers = new ArucoCorners;
    markers->corners = new ArucoCorner[markerCorners.size()];
    markers->length = markerCorners.size();

    for (size_t i = 0; i < markerCorners.size(); ++i) {
        ArucoCorner c = {
            markerIds[i],
            Point2f{markerCorners[i][0].x, markerCorners[i][0].y},
            Point2f{markerCorners[i][1].x, markerCorners[i][1].y},
            Point2f{markerCorners[i][2].x, markerCorners[i][2].y},
            Point2f{markerCorners[i][3].x, markerCorners[i][3].y},
        };
        markers->corners[i] = c;
    }

    ArucoCorners* rejected = new ArucoCorners;
    rejected->corners = new ArucoCorner[markerRejected.size()];
    rejected->length = markerRejected.size();
    for (size_t i = 0; i < markerRejected.size(); ++i) {
        ArucoCorner c = {
            -1,
            Point2f{markerRejected[i][0].x, markerRejected[i][0].y},
            Point2f{markerRejected[i][1].x, markerRejected[i][1].y},
            Point2f{markerRejected[i][2].x, markerRejected[i][2].y},
            Point2f{markerRejected[i][3].x, markerRejected[i][3].y},
        };
        rejected->corners[i] = c;
    }

    ArucoDetectedMarkers ret = {
        markers,
        rejected
    };

    return ret;
}

void ArucoDrawAxis(Mat dst, Mat cameraMatrix, Mat distCoeffs, Mat rvec, Mat tvec, float length){
    cv::aruco::drawAxis(*dst, *cameraMatrix, *distCoeffs, *rvec, *tvec, length);
}

void ArucoCorners_Close(ArucoCorners* c) {
    delete c;
}

void ArucoEstimatePoseSingleMarkers(ArucoCorners* corners, float markerLength, Mat cameraMatrix, Mat distCoeffs, struct Mats* outRvec, struct Mats* outTvec) {
    std::vector<cv::Vec3d> rvecs, tvecs;
    std::vector<std::vector<cv::Point2f>> cornersVec;

    for (size_t i = 0; i < corners->length; ++i) {
        std::vector<cv::Point2f> cc;
        cc.push_back(cv::Point2f(corners->corners[i].p1.x, corners->corners[i].p1.y));
        cc.push_back(cv::Point2f(corners->corners[i].p2.x, corners->corners[i].p2.y));
        cc.push_back(cv::Point2f(corners->corners[i].p3.x, corners->corners[i].p3.y));
        cc.push_back(cv::Point2f(corners->corners[i].p4.x, corners->corners[i].p4.y));
        cornersVec.push_back(cc);
    }
    cv::aruco::estimatePoseSingleMarkers(cornersVec, markerLength, *cameraMatrix, *distCoeffs, rvecs, tvecs);

    outRvec->mats = new Mat[rvecs.size()];
    outRvec->length = rvecs.size();
    for (size_t i = 0; i < rvecs.size(); ++i) {
        cv::Mat* m = new cv::Mat(3,1, CV_64F);
        m->at<double>(0,0) = rvecs[i][0];
        m->at<double>(1,0) = rvecs[i][1];
        m->at<double>(2,0) = rvecs[i][2];
        outRvec->mats[i] = m;
    }
    outTvec->length = rvecs.size();
    outTvec->mats = new Mat[rvecs.size()];
    for (size_t i = 0; i < tvecs.size(); ++i) {
        cv::Mat* m = new cv::Mat(3,1, CV_64F);
        m->at<double>(0,0) = tvecs[i][0];
        m->at<double>(1,0) = tvecs[i][1];
        m->at<double>(2,0) = tvecs[i][2];
        outTvec->mats[i] = m;
    }
}
