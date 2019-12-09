package file

import "github.com/aws/aws-sdk-go/service/s3"

type AccessControlList string
type ContentType string

const (
	AccessControlPrivate         AccessControlList = "private"
	AccessControlPublicRead      AccessControlList = "public-read"
	AccessControlPublicReadWrite AccessControlList = "public-read-write"
	ContentTypeJpeg              ContentType       = "image/jpeg"
)

func (a AccessControlList) toAWSACL() *string {
	switch a {
	case AccessControlPrivate:
		return strToPtr(s3.ObjectCannedACLPrivate)
	case AccessControlPublicRead:
		return strToPtr(s3.ObjectCannedACLPublicRead)
	case AccessControlPublicReadWrite:
		return strToPtr(s3.ObjectCannedACLPublicReadWrite)
	default:
		return nil
	}
}

func (c ContentType) toS3ContentType() *string {
	return strToPtr(string(c))
}
