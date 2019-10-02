package file

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/snapmartinc/logger"
)

type CSVProcessor interface {
	DownloadFile(ctx context.Context, fileName, filePath string) (Result, error)
	UploadResult(ctx context.Context, result Result) error
	RemoveFiles(ctx context.Context, result Result)
}

func NewCSVFileHandler(
	loggerFactory logger.Factory,
	fileManager Manager,
) CSVProcessor {
	return &csvFileHandler{
		loggerFactory: loggerFactory,
		fileManager:   fileManager,
	}
}

type csvFileHandler struct {
	loggerFactory logger.Factory
	fileManager   Manager
}

type Result struct {
	InputFile  InputFile
	OutputFile OutputFile
}

type InputFile struct {
	Body io.Reader
	Name string
	Path string
}

type OutputFile struct {
	Body io.ReadWriteSeeker
	Name string
	Path string
}

// Handle handle journey status update request
func (c *csvFileHandler) DownloadFile(ctx context.Context, fileName, filePath string) (Result, error) {
	c.loggerFactory.
		Logger(ctx).
		WithField("csvProcessor", "Download and create temporary files").
		Info(fmt.Sprintf("fileName:%s, filePath:%s", fileName, filePath))

	uploadedFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create temporary file to hold uploaded", err)
	}
	reportFile, err := os.Create(fmt.Sprintf("report-%s", fileName))
	if err != nil {
		log.Fatal("Cannot create temporary file to hold report", err)
	}
	result := Result{
		InputFile: InputFile{
			Body: uploadedFile,
			Name: uploadedFile.Name(),
			Path: filePath,
		},
		OutputFile: OutputFile{
			Body: reportFile,
			Name: reportFile.Name(),
			Path: filePath,
		},
	}
	if err := c.fileManager.DownloadFile(ctx, uploadedFile, filePath, fileName); err != nil {
		return result, err
	}
	return result, nil
}

// Handle handle journey status update request
func (c *csvFileHandler) UploadResult(ctx context.Context, result Result) error {
	c.loggerFactory.
		Logger(ctx).
		WithField("csvProcessor", "Upload result file").
		Info(fmt.Sprintf("upload report file"))

	temp, _ := os.Open(result.OutputFile.Name)
	defer func() {
		err := temp.Close()
		if err != nil {
			c.loggerFactory.Logger(ctx).WithField("csvProcessor", "Upload report file").Error(err)
		}
	}()
	if err := c.fileManager.UploadFiles(ctx, []File{
		{
			Path: result.OutputFile.Path,
			Name: result.OutputFile.Name,
			Body: temp,
		},
	}); err != nil {
		return err
	}
	return nil
}

func (c *csvFileHandler) RemoveFiles(ctx context.Context, result Result) {
	c.loggerFactory.
		Logger(ctx).
		WithField("csvProcessor", "remove temporary files").
		Info(fmt.Sprintf("remove files: %s, %s", result.InputFile.Name, result.OutputFile.Name))

	if err := os.Remove(result.InputFile.Name); err != nil {
		log.Fatalf("Cannot create delete file %s, error: %s", result.InputFile.Name, err)
	}
	if err := os.Remove(result.OutputFile.Name); err != nil {
		log.Fatalf("Cannot create delete file %s, error: %s", result.OutputFile.Name, err)
	}
}
