package file

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type Manager interface {
	UploadFiles(ctx context.Context, files []File) error
	DownloadFile(ctx context.Context, destination io.WriterAt, path, name string) error
}

type s3Manager struct {
	Bucket string
	Cfg    *aws.Config
}

type S3Config struct {
	Region       string `envconfig:"REGION" required:"true"`
	Bucket       string `envconfig:"BUCKET" required:"true"`
	AwsAccessKey string `envconfig:"ACCESS_KEY" required:"true"`
	AwsSecret    string `envconfig:"SECRET" required:"true"`
}

func NewS3Manager(conf S3Config) (Manager, error) {
	creds := credentials.NewStaticCredentials(conf.AwsAccessKey, conf.AwsSecret, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	cfg := aws.NewConfig().WithRegion(conf.Region).WithCredentials(creds)
	return s3Manager{Bucket: conf.Bucket, Cfg: cfg}, nil
}

type File struct {
	Path        string
	Name        string
	Body        io.Reader
	ACL         AccessControlList
	ContentType ContentType
}

func (u s3Manager) UploadFiles(ctx context.Context, files []File) error {
	sess, err := session.NewSession(u.Cfg)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	for i := range files {
		_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			ACL:         files[i].ACL.toAWSACL(),
			Bucket:      aws.String(u.Bucket),
			Body:        files[i].Body,
			Key:         aws.String(getLocation(files[i].Path, files[i].Name)),
			ContentType: files[i].ContentType.toS3ContentType(),
		})
		if err != nil {
			return errors.Wrap(err, "Cannot upload file: "+files[i].Path)
		}
	}
	return nil
}

func (u s3Manager) DownloadFile(ctx context.Context, destination io.WriterAt, path, name string) error {
	sess, err := session.NewSession(u.Cfg)
	if err != nil {
		return err
	}
	uploader := s3manager.NewDownloader(sess)
	_, err = uploader.Download(destination, &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(getLocation(path, name)),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot download file: %s", getLocation(path, name)))
	}
	return nil
}

func getLocation(path, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}
