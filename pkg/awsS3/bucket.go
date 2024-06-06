package awsS3

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const (
	awsRegionEnv = "AWS_REGION"
	awsKeyEnv    = "AWS_ACCESS_KEY_ID"
	awsSecretEnv = "AWS_SECRET_ACCESS_KEY"
	maxImageSize = 5 * 1024 * 1024 // 5 MB
)

type S3ServiceI interface {
	UploadPhotos(photos []string, photoNamesPrefix string) ([]models.Photo, error)
	UploadSinglePhoto(photoURI string, prefix string) (*manager.UploadOutput, error)
}

type S3Service struct {
	uploader *manager.Uploader
	bucket   string
}

func NewS3Service(bucket string) *S3Service {
	client := s3.New(s3.Options{
		Region: viper.GetString(awsRegionEnv),
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				viper.GetString(awsKeyEnv),
				viper.GetString(awsSecretEnv), "",
			)),
	})

	return &S3Service{
		uploader: manager.NewUploader(client),
		bucket:   bucket,
	}
}

// UploadSinglePhoto uploads a single photo to S3.
func (s *S3Service) UploadSinglePhoto(photoURI string, prefix string) (*manager.UploadOutput, error) {
	valid, decodedData, err := s.validateAndDecodeBase64(photoURI)
	if err != nil || !valid {
		return nil, err
	}

	result, err := s.uploadImage(bytes.NewReader(decodedData), s.createFilename(prefix))
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UploadPhotos uploads multiple photos to S3.
func (s *S3Service) UploadPhotos(photos []string, photoNamesPrefix string) ([]models.Photo, error) {
	photoUrls := make([]models.Photo, 0, len(photos))
	for _, photoURI := range photos {
		result, err := s.UploadSinglePhoto(photoURI, photoNamesPrefix)
		if err != nil {
			return nil, err
		}

		photoUrls = append(photoUrls, models.Photo{ImageURL: result.Location, Key: *result.Key})
	}
	return photoUrls, nil
}

func (s *S3Service) uploadImage(file io.Reader, filename string) (*manager.UploadOutput, error) {
	result, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename + ".jpg"),
		Body:   file,
		ACL:    "public-read",
	})

	return result, err
}

func (s *S3Service) validateAndDecodeBase64(photoURI string) (bool, []byte, error) {
	data, err := extractBase64Part(photoURI)
	if err != nil {
		return false, nil, err
	}
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return false, nil, err
	}
	if len(decodedData) > maxImageSize {
		return false, nil, errors.New("image is too large")
	}
	return true, decodedData, nil
}

func extractBase64Part(dataURI string) (string, error) {
	parts := strings.Split(dataURI, ",")
	if len(parts) != 2 {
		return "", errors.New("incorrect dataURI format")
	}
	return parts[1], nil
}

func (s *S3Service) createFilename(namePrefix string) string {
	id := uuid.New()
	result := fmt.Sprintf("%s_%s", namePrefix, id)
	return result
}
