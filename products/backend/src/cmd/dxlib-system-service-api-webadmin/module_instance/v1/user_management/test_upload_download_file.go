package user_management

import (
	"bytes"
	"fmt"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/object_storage"
	utilsHttp "github.com/donnyhardyanto/dxlib/utils/http"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
)

func testUploadFile(aepr *api.DXAPIEndPointRequest) (err error) {
	const maxRequestSize = 100 * 1024 * 1024 // 100MB

	// Check the request size
	if aepr.Request.ContentLength > maxRequestSize {
		return aepr.WriteResponseAndNewErrorf(http.StatusRequestEntityTooLarge, "", "REQUEST_ENTITY_TOO_LARGE")
	}

	_, filename, err := aepr.GetParameterValueAsString("filename")
	if err != nil {
		return err
	}
	_, objectStorageNameId, err := aepr.GetParameterValueAsString("object_storage_nameid")
	if err != nil {
		return err
	}

	err = object_storage.Manager.FindObjectStorageAndReceiveObject(aepr, objectStorageNameId, filename)
	if err != nil {
		return err
	}

	return nil
}

func testUploadFile2(aepr *api.DXAPIEndPointRequest) (err error) {
	const maxRequestSize = 100 * 1024 * 1024 // 100MB

	// Check the request size
	if aepr.Request.ContentLength > maxRequestSize {
		return aepr.WriteResponseAndNewErrorf(http.StatusRequestEntityTooLarge, "", "REQUEST_ENTITY_TOO_LARGE")
	}

	_, filename, err := aepr.GetParameterValueAsString("filename")
	if err != nil {
		return err
	}
	_, objectStorageNameId, err := aepr.GetParameterValueAsString("object_storage_nameid")
	if err != nil {
		return err
	}

	objectStorage, exists := object_storage.Manager.ObjectStorages[objectStorageNameId]
	if !exists {
		return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "OBJECT_STORAGE_NAME_NOT_FOUND:%s", objectStorageNameId)
	}

	bodyLen := aepr.Request.ContentLength
	aepr.Log.Infof("Request body length: %d", bodyLen)

	bs := aepr.Request.Body
	if bs == nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "BODY_IS_NIL")
	}

	// RequestRead the entire request body into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, bs)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "REQUEST_BODY_READ_FAILED:%v", err.Error())
	}

	// Upload the original file

	aBuf := buf.Bytes()
	aBufLen := int64(len(aBuf))

	uploadInfo, err := objectStorage.UploadStream(bytes.NewReader(aBuf), filename, filename, "application/octet-stream", false, aBufLen)
	if err != nil {
		return err
	}

	aepr.Log.Infof("Upload info result: %v", uploadInfo)

	// Decode the image
	img, formatName, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusInternalServerError, "", "DECODE_IMAGE_FAILED:%v", err.Error())
	}

	targetWidth := 128
	targetHeight := 128
	quality := 100

	aepr.Log.Infof("Image Request format name: %s", formatName)

	formatName = "png"
	resizedImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf2 bytes.Buffer
	switch formatName {
	case "jpeg":
		err = jpeg.Encode(&buf2, resizedImg, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf2, resizedImg)
	default:
		return errors.Errorf("UNSUPPORTED_FORMAT: %s", formatName)
	}
	if err != nil {
		return err
	}

	filename2 := filename + ".resized." + formatName

	aBufMod := buf2.Bytes()
	aBufModLen := int64(len(aBufMod))

	uploadInfo2, err := objectStorage.UploadStream(bytes.NewReader(aBufMod), filename2, filename2, "image/"+formatName, false, aBufModLen)

	//	uploadInfo2, err := objectStorage.UploadStream(&buf2, filename2, filename2, "application/octet-stream")
	if err != nil {
		return err
	}

	aepr.Log.Infof("Upload 2  info result: %v", uploadInfo2)
	return nil
}

func testUploadFile3(aepr *api.DXAPIEndPointRequest) (err error) {
	const maxRequestSize = 100 * 1024 * 1024 // 100MB

	// Check the request size
	if aepr.Request.ContentLength > maxRequestSize {
		return aepr.WriteResponseAndNewErrorf(http.StatusRequestEntityTooLarge, "", "REQUEST_ENTITY_TOO_LARGE")
	}

	_, filename, err := aepr.GetParameterValueAsString("filename")
	if err != nil {
		return err
	}
	_, objectStorageNameId, err := aepr.GetParameterValueAsString("object_storage_nameid")
	if err != nil {
		return err
	}

	objectStorage, exists := object_storage.Manager.ObjectStorages[objectStorageNameId]
	if !exists {
		return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "OBJECT_STORAGE_NAME_NOT_FOUND:%s", objectStorageNameId)
	}

	bodyLen := aepr.Request.ContentLength
	aepr.Log.Infof("Request body length: %d", bodyLen)

	// Get the request body stream
	bs := aepr.Request.Body
	if bs == nil {
		return errors.Errorf("FAILED_TO_GET_BODY_STREAM: body is nil")
	}

	// RequestRead the entire request body into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, bs)
	if err != nil {
		return errors.Errorf("failed to read request body: %v", err.Error())
	}

	// Upload the original file
	aBuf := buf.Bytes()
	aBufLen := int64(len(aBuf))

	uploadInfo, err := objectStorage.UploadStream(bytes.NewReader(aBuf), filename, filename, "application/octet-stream", false, aBufLen)
	if err != nil {
		return errors.Errorf("failed to upload original file: %v", err.Error())
	}

	aepr.Log.Infof("Original upload info result: %v", uploadInfo)

	// Decode the image
	img, formatName, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return errors.Errorf("failed to decode image: %v", err.Error())
	}

	aepr.Log.Infof("Image format: %s", formatName)

	// Define resize dimensions
	sizes := []struct {
		width  int
		height int
	}{
		{128, 128},
		{64, 64},
		{32, 32},
	}

	quality := 100

	for _, size := range sizes {
		// Resize the image
		resizedImg := image.NewRGBA(image.Rect(0, 0, size.width, size.height))
		draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

		// Encode the resized image
		var resizedBuf bytes.Buffer
		switch formatName {
		case "jpeg":
			err = jpeg.Encode(&resizedBuf, resizedImg, &jpeg.Options{Quality: quality})
		case "png":
			err = png.Encode(&resizedBuf, resizedImg)
		default:
			return errors.Errorf("UNSUPPORTED_FORMAT: %s", formatName)
		}
		if err != nil {
			return errors.Errorf("failed to encode resized image (%dx%d): %v", size.width, size.height, err.Error())
		}

		// Upload the resized image
		resizedFilename := fmt.Sprintf("%s.resized.%dx%d.%s", filename, size.width, size.height, formatName)
		aBuf := resizedBuf.Bytes()
		aBufLen := int64(len(aBuf))

		uploadInfo, err := objectStorage.UploadStream(bytes.NewReader(aBuf), resizedFilename, resizedFilename, "image/"+formatName, false, aBufLen)
		if err != nil {
			return errors.Errorf("failed to upload resized image (%dx%d): %v", size.width, size.height, err.Error())
		}

		aepr.Log.Infof("Resized (%dx%d) upload info result: %v", size.width, size.height, uploadInfo)
	}

	return nil
}

func testDownloadFile(aepr *api.DXAPIEndPointRequest) (err error) {
	_, filename, err := aepr.GetParameterValueAsString("filename")
	if err != nil {
		return err
	}

	_, objectStorageNameId, err := aepr.GetParameterValueAsString("object_storage_nameid")
	if err != nil {
		return err
	}

	err = object_storage.Manager.FindObjectStorageAndSendObject(aepr, objectStorageNameId, filename)
	if err != nil {
		return err
	}

	return nil
}

func defineAPITestUploadDownloadFile(anAPI *api.DXAPI) {
	anAPI.NewEndPoint("TestUploadFile",
		"Test Upload File",
		"/v1/test/upload/file", "POST", api.EndPointTypeHTTPUploadStream, utilsHttp.ContentTypeApplicationOctetStream, []api.DXAPIEndPointParameter{
			{NameId: "filename", Type: "string", Description: "Filename", IsMustExist: true},
			{NameId: "object_storage_nameid", Type: "string", Description: "Object Storage NameId", IsMustExist: true},
		}, testUploadFile3, nil, nil, nil, nil, 0, "default",
	)

	anAPI.NewEndPoint("TestDownloadFile",
		"Test Download File",
		"/v1/test/download/file", "POST", api.EndPointTypeHTTPDownloadStream, utilsHttp.ContentTypeApplicationJSON, []api.DXAPIEndPointParameter{
			{NameId: "filename", Type: "string", Description: "Filename", IsMustExist: true},
			{NameId: "object_storage_nameid", Type: "string", Description: "Object Storage NameId", IsMustExist: true},
		}, testDownloadFile, nil, nil, nil, nil, 0, "default",
	)

}
