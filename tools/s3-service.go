package tools

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oklog/ulid/v2"
	"github.com/vincent-petithory/dataurl"
	"github.com/zihaolam/roadwarrior-backend/config"
	"github.com/zihaolam/roadwarrior-backend/database"
)

type S3Object struct {
	Bucket string
	Key    string
}

var s3Client = s3.New(database.AWSSession, &aws.Config{Region: aws.String("ap-southeast-1")})

func UploadBase64FileToS3(base64File string) (*S3Object, error) {
	dataURL, err := dataurl.DecodeString(base64File)
	if err != nil {
		return nil, err
	}

	fileName := ulid.Make().String() + "." + dataURL.Subtype
	params := &s3.PutObjectInput{Bucket: aws.String(config.BucketName), Key: aws.String(fileName), Body: aws.ReadSeekCloser(strings.NewReader(string(dataURL.Data)))}
	_, err = s3Client.PutObject(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &S3Object{Bucket: config.BucketName, Key: fileName}, nil
}

func GetS3ObjectUrl(s3Object *S3Object, region string) string {
	return fmt.Sprintf("https://%s.amazonaws.com/%s/%s", region, s3Object.Bucket, s3Object.Key)
}

func GetBucketKeyFromS3ObjectUrl(s3ObjectUrl string) (string, string, error) {
	u, err := url.Parse(s3ObjectUrl)
	if err != nil {
		return "", "", err
	}

	// The path should contain the bucket and key information
	path := strings.TrimPrefix(u.Path, "/") // Remove leading slash if present
	parts := strings.SplitN(path, "/", 2)   // Split into two parts: bucket and key

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid S3 object URL format")
	}

	return parts[0], parts[1], nil
}

func DeleteObjectFromS3(bucket string, key string) error {
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	return err
}
