package lib

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/donnyhardyanto/dxlib/api"
	"github.com/donnyhardyanto/dxlib/object_storage"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"time"
)

type ProcessedImageObjectStorage struct {
	ObjectStorageNameId string
	Width               int
	Height              int
}

type ImageObjectStorage struct {
	ObjectStorageSourceNameId string
	MaxRequestSize            int64
	MaxPixelWidth             int64
	MaxPixelHeight            int64
	MaxBytesPerPixel          int64
	MaxPixels                 int64
	ProcessedImages           map[string]ProcessedImageObjectStorage
}

const MaxRequestSize = 100 * 1024 * 1024 // 100MB
const MaxPixelWidth = 4096
const MaxPixelHeight = 4096
const MaxBytesPerPixel = 10
const MaxPixels = 40000000                // ~40MP
const ImageProcessLimit = 5 * time.Second // Timeout for image processing

func NewImageObjectStorage(objectStorageSourceNameId string,
	maxRequestSize int64, maxPixelWidth int64, maxPixelHeight int64, maxBytesPerPixel int64, maxPixels int64,
	processedImages map[string]ProcessedImageObjectStorage) *ImageObjectStorage {
	return &ImageObjectStorage{
		ObjectStorageSourceNameId: objectStorageSourceNameId,
		MaxRequestSize:            maxRequestSize,
		MaxPixelWidth:             maxPixelWidth,
		MaxPixelHeight:            maxPixelHeight,
		MaxBytesPerPixel:          maxBytesPerPixel,
		MaxPixels:                 maxPixels,
		ProcessedImages:           processedImages,
	}
}

func calculateAspectRatioHeight(originalWidth, originalHeight, targetWidth int) int {
	ratio := float64(originalHeight) / float64(originalWidth)
	return int(float64(targetWidth) * ratio)
}

// ValidateImageDimensions checks for pixel flood attacks by validating image dimensions
func (ios *ImageObjectStorage) ValidateImageDimensions(data []byte) error {
	// Create a context with timeout for image processing
	ctx, cancel := context.WithTimeout(context.Background(), ImageProcessLimit)
	defer cancel()

	resultCh := make(chan struct {
		cfg    image.Config
		format string
		err    error
	}, 1)

	// Use a goroutine to process the image with a timeout
	go func() {
		reader := bytes.NewReader(data)
		cfg, format, err := image.DecodeConfig(reader)
		resultCh <- struct {
			cfg    image.Config
			format string
			err    error
		}{cfg, format, err}
	}()

	select {
	case <-ctx.Done():
		return errors.New("IMAGE_PROCESSING_TIMEOUT:possible_decompression_bomb")
	case result := <-resultCh:
		if result.err != nil {
			return errors.Errorf("FAILED_TO_DECODE_IMAGE_CONFIG:%v", result.err)
		}

		// Validate image format
		if result.format == "" {
			return errors.New("UNKNOWN_IMAGE_FORMAT")
		}

		// Check dimensions against limits
		if result.cfg.Width <= 0 || result.cfg.Height <= 0 {
			return errors.New("INVALID_IMAGE_DIMENSIONS:dimensions_cannot_be_zero_or_negative")
		}

		if int64(result.cfg.Width) > ios.MaxPixelWidth {
			return errors.Errorf("IMAGE_WIDTH_EXCEEDS_LIMIT:max=%d", ios.MaxPixelWidth)
		}

		if int64(result.cfg.Height) > ios.MaxPixelHeight {
			return errors.Errorf("IMAGE_HEIGHT_EXCEEDS_LIMIT:max=%d", ios.MaxPixelHeight)
		}

		// Check total pixels
		totalPixels := int64(result.cfg.Width) * int64(result.cfg.Height)
		if totalPixels > ios.MaxPixels {
			return errors.Errorf("TOTAL_PIXELS_EXCEEDS_LIMIT:max=%d", ios.MaxPixels)
		}

		// Check pixels-per-byte ratio to detect compression bombs
		fileSize := int64(len(data))
		if fileSize == 0 {
			return errors.New("INVALID_FILE_SIZE:size_is_zero")
		}

		pixelsPerByte := float64(totalPixels) / float64(fileSize)
		if pixelsPerByte > float64(ios.MaxBytesPerPixel) {
			return errors.Errorf("SUSPICIOUS_PIXELS_TO_FILESIZE_RATIO:ratio=%.2f", pixelsPerByte)
		}

		return nil
	}
}

// DecodeImageWithTimeout decodes an image with a timeout to prevent DoS attacks
func (ios *ImageObjectStorage) DecodeImageWithTimeout(data []byte) (image.Image, string, error) {
	// Create a context with timeout for image processing
	ctx, cancel := context.WithTimeout(context.Background(), ImageProcessLimit)
	defer cancel()

	resultCh := make(chan struct {
		img    image.Image
		format string
		err    error
	}, 1)

	// Use a goroutine to process the image with a timeout
	go func() {
		img, format, err := image.Decode(bytes.NewReader(data))
		resultCh <- struct {
			img    image.Image
			format string
			err    error
		}{img, format, err}
	}()

	// Wait for either the result or timeout
	select {
	case <-ctx.Done():
		return nil, "", errors.New("IMAGE_DECODE_TIMEOUT:possible_decompression_bomb")
	case result := <-resultCh:
		return result.img, result.format, result.err
	}
}

func (ios *ImageObjectStorage) Update(aepr *api.DXAPIEndPointRequest, filename string, fileContentBase64 string) (err error) {

	// Check the request size
	if aepr.Request.ContentLength > ios.MaxRequestSize {
		return aepr.WriteResponseAndNewErrorf(http.StatusRequestEntityTooLarge, "", "REQUEST_ENTITY_TOO_LARGE")
	}

	objectStorage, exists := object_storage.Manager.ObjectStorages[ios.ObjectStorageSourceNameId]
	if !exists {
		return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "OBJECT_STORAGE_NAME_NOT_FOUND:%s", ios.ObjectStorageSourceNameId)
	}

	var buf bytes.Buffer
	var bodyLen int64
	if fileContentBase64 == "" {
		bodyLen = aepr.Request.ContentLength
		aepr.Log.Infof("Request body length: %d", bodyLen)

		// Get the request body stream
		bs := aepr.Request.Body
		if bs == nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_GET_BODY_STREAM:%s", ios.ObjectStorageSourceNameId)
		}

		// RequestRead the entire request body into a buffer
		_, err = io.Copy(&buf, bs)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_READ_REQUEST_BODY:%s=%v", ios.ObjectStorageSourceNameId, err.Error())
		}
	} else {
		// Decode base64 string to bytes
		decodedBytes, err := base64.StdEncoding.DecodeString(fileContentBase64)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_DECODE_BASE64:%s=%v", ios.ObjectStorageSourceNameId, err.Error())
		}

		// Get the total size of the decoded content
		bodyLen = int64(len(decodedBytes))
		aepr.Log.Infof("Base64 decoded content length: %d", bodyLen)

		// Write decoded bytes to buffer
		buf.Write(decodedBytes)
	}

	// Validate image dimensions to prevent pixel flood attacks
	err = ios.ValidateImageDimensions(buf.Bytes())
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "PIXEL_FLOOD_PROTECTION:%v", err)
	}

	// Upload the original file
	uploadInfo, err := objectStorage.UploadStream(bytes.NewReader(buf.Bytes()), filename, filename, "application/octet-stream", false, bodyLen)
	if err != nil {
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_UPLOAD_SOURCE_IMAGE_TO_OBJECT_STORAGE:%s=%v", ios.ObjectStorageSourceNameId, err.Error())
	}

	aepr.Log.Infof("Original upload info result: %d", uploadInfo.Size)

	// Decode the image with timeout protection
	img, formatName, err := ios.DecodeImageWithTimeout(buf.Bytes())
	if err != nil {
		if err.Error() == "IMAGE_DECODE_TIMEOUT:possible_decompression_bomb" {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "PIXEL_FLOOD_PROTECTION:%v", err)
		}
		return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_DECODE_IMAGE:%s=%v", ios.ObjectStorageSourceNameId, err)
	}

	aepr.Log.Infof("Image format (using Image.Decode): %s", formatName)

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	for _, processedImage := range ios.ProcessedImages {
		objectStorage, ok := object_storage.Manager.ObjectStorages[processedImage.ObjectStorageNameId]
		if !ok {
			return aepr.WriteResponseAndNewErrorf(http.StatusNotFound, "", "OBJECT_STORAGE_NAME_NOT_FOUND:%s", processedImage.ObjectStorageNameId)
		}

		// Set a maximum target height based on aspect ratio and configured width
		targetHeight := calculateAspectRatioHeight(originalWidth, originalHeight, processedImage.Width)

		// Safeguard against extremely tall images
		if int64(targetHeight) > ios.MaxPixelHeight {
			targetHeight = int(ios.MaxPixelHeight)
		}

		// Create a context with timeout for image scaling
		ctx, cancel := context.WithTimeout(context.Background(), ImageProcessLimit)
		resizedImg := image.NewRGBA(image.Rect(0, 0, processedImage.Width, targetHeight))

		// Use a goroutine to scale the image with a timeout
		scaleDone := make(chan bool, 1)
		var scaleErr error

		go func() {
			// Resize the image
			draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
			scaleDone <- true
		}()

		// Wait for either completion or timeout
		select {
		case <-ctx.Done():
			cancel()
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "PIXEL_FLOOD_PROTECTION:image_scaling_timeout")
		case <-scaleDone:
			cancel()
			if scaleErr != nil {
				return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_SCALE_IMAGE:%v", scaleErr.Error())
			}
		}

		// Encode the resized image
		var resizedBuf bytes.Buffer
		err = png.Encode(&resizedBuf, resizedImg)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "RESIZED_IMAGE_PNG_ENCODE_FAILED:(%dx%d) %v", processedImage.Width, targetHeight, err.Error())
		}

		// Upload the resized image
		buf := resizedBuf.Bytes()
		bufLen := int64(len(buf))
		uploadInfo, err := objectStorage.UploadStream(bytes.NewReader(buf), filename, filename, "image/"+formatName, false, bufLen)
		if err != nil {
			return aepr.WriteResponseAndNewErrorf(http.StatusUnprocessableEntity, "", "FAILED_TO_UPLOAD_RESIZED_IMAGE_TO_OBJECT_STORAGE:(%s)=%v", processedImage.ObjectStorageNameId, err.Error())
		}

		aepr.Log.Infof("Resized (%dx%d) upload info result size: %d", processedImage.Width, targetHeight, uploadInfo.Size)
	}

	return nil
}

func (ios *ImageObjectStorage) DownloadSource(aepr *api.DXAPIEndPointRequest, filename string) (err error) {
	err = object_storage.Manager.FindObjectStorageAndSendObject(aepr, ios.ObjectStorageSourceNameId, filename)
	if err != nil {
		return err
	}

	return nil
}

func (ios *ImageObjectStorage) DownloadProcessedImage(aepr *api.DXAPIEndPointRequest, processedImageNameId string, filename string) (err error) {

	err = object_storage.Manager.FindObjectStorageAndSendObject(aepr, ios.ProcessedImages[processedImageNameId].ObjectStorageNameId, filename)
	if err != nil {
		return err
	}

	return nil
}
