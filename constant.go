package file

import "github.com/aws/aws-sdk-go/service/s3"

type AccessControlList string

const (
	AccessControlPrivate         AccessControlList = "private"
	AccessControlPublicRead      AccessControlList = "public"
	AccessControlPublicReadWrite AccessControlList = "public-read-write"
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
