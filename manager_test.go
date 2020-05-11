package file_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	upload_files "github.com/best-expendables/upload-files"
)

func TestS3Uploader_Upload(t *testing.T) {
	t.SkipNow()
	conf := upload_files.S3Config{
		Region:       "ap-southeast-1",
		Bucket:       "snapmart-staging",
		AwsAccessKey: "AKIATPQMCETQLYEWVGHF",
		AwsSecret:    "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",
	}
	uploader, err := upload_files.NewS3Manager(conf)
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	f, err := os.Open("temp.txt")
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}

	err = uploader.UploadFiles(context.Background(), []upload_files.File{{
		Path: "new/test/file",
		Name: "text.txt",
		Body: f,
		ACL:  upload_files.AccessControlPublicRead,
	}})
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
}

func TestS3Downloader(t *testing.T) {
	t.SkipNow()
	conf := upload_files.S3Config{
		Region:       "ap-southeast-1",
		Bucket:       "snapmart-staging",
		AwsAccessKey: "AKIATPQMCETQLYEWVGHF",
		AwsSecret:    "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",
	}
	uploader, err := upload_files.NewS3Manager(conf)
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	fileName := "temp.txt"
	f, err := os.Open(fileName)
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}
	err = uploader.UploadFiles(context.Background(), []upload_files.File{{
		Path: "new/test/file",
		Name: fileName,
		Body: f,
		ACL:  upload_files.AccessControlPublicRead,
	}})
	if err != nil {
		t.Errorf("error not nil: %v", err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()
	if err := uploader.DownloadFile(context.Background(), tmpFile, "new/test/file", "temp.txt"); err != nil {
		t.Errorf("error not nil: %v", err)
	}
	bs := csv.NewReader(tmpFile)
	for {
		line, err := bs.Read()
		fmt.Println(line, err)
		if err == io.EOF {
			break
		}
	}

}
