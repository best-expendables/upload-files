package file_test

import (
	"os"
	upload_files "snapmartinc/upload-files"
	"testing"
)

func TestS3Uploader_Upload(t *testing.T) {
	t.SkipNow()
	conf := upload_files.S3Config{
		Region:       "ap-southeast-1",
		Bucket:       "quangphan-test",
		AwsAccessKey: "AKIAR75ZZDHPW42OWS4X",
		AwsSecret:    "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",
	}
	uploader, err := upload_files.NewS3Manager(conf)
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	f, err := os.Open("uploader.go")
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	err = uploader.Upload("new/test/file", f)
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
}
