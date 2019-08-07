package file

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type FileManager interface {
	Upload(path string, body io.Reader) error
}

type s3Manager struct {
	Bucket string
	Cfg    *aws.Config
}

type S3Config struct {
	Region       string
	Bucket       string
	AwsAccessKey string
	AwsSecret    string
}

func NewS3Manager(conf S3Config) (FileManager, error) {
	creds := credentials.NewStaticCredentials(conf.AwsAccessKey, conf.AwsSecret, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	cfg := aws.NewConfig().WithRegion(conf.Region).WithCredentials(creds)
	return s3Manager{Bucket: conf.Bucket, Cfg: cfg}, nil
}

func (u s3Manager) Upload(path string, body io.Reader) error {
	sess, err := session.NewSession(u.Cfg)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.Bucket),
		Body:   body,
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}
