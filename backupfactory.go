package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// This function does the actual backup
func BackupFactory(config Config) {

	// Create a new aws session using credentials
	// from the specified aws credentials profile
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-2"),
		Credentials: credentials.NewSharedCredentials("", config.AwsProfile),
	})

	if err != nil {
		exitErrorf("Error creating AWS Session %v", err)
	}

	// Create a new instance of an s3manager using the
	// above session
	uploader := s3manager.NewUploader(sess)

	var filesToUpload []string

	// using filepath.Walk to create a list of files we
	// want to upload to s3
	// This basically does a recursive scan of the paths we
	// have specified in items.st / config.Directories
	for i := 0; i < len(config.Directories); i++ {
		err = filepath.Walk(config.Directories[i], func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				filesToUpload = append(filesToUpload, path)
			}
			return nil
		})
		if err != nil {
			exitErrorf("Error", err)
		}
	}

	// Loop over the length of filesToUpload
	// and send each file to aws s3
	for i := 0; i < len(filesToUpload); i++ {
		filename := filesToUpload[i]
		bucket := config.TargetBucket

		file, err := os.Open(filename)
		if err != nil {
			exitErrorf("Unable to open file %q, %v", err)
		}

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   file,
		})
		if err != nil {
			exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
		}

		// Woop to the user
		fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)

	}

}
