package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Handler(ctx context.Context, S3Event events.S3Event) {

	// create file in /tmp/ to write data from bucket
	file, createerr := os.Create("/tmp/" + S3Event.Records[0].S3.Object.Key)
	if createerr != nil {
		exitErrorf("Unable to open file %q, %v", S3Event.Records[0].S3.Object.Key, createerr)
	}

	defer file.Close()

	// creating new needed session
	mySession := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(mySession)

	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(S3Event.Records[0].S3.Bucket.Name),
			Key:    aws.String(S3Event.Records[0].S3.Object.Key),
		})
	if err != nil {
		exitErrorf("Unable to download s3File %q, %v", S3Event.Records[0].S3.Object.Key, err)
	}

	dat, readerr := ioutil.ReadFile(file.Name())

	if readerr != nil {
		exitErrorf("Cannot read the file", readerr)
	}

	//ToDo: implement transformation

	json_filename := strings.Replace(S3Event.Records[0].S3.Object.Key, ".xml", ".json", 1)

	_ = ioutil.WriteFile("/tmp/"+json_filename, dat, 0644)

	json_file, err := os.Open("/tmp/" + json_filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer json_file.Close()

	// Uploading to S3
	uploader := s3manager.NewUploader(mySession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("OUTPUT_BUCKET")),

		Key: aws.String(json_filename),

		Body: json_file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", json_filename, os.Getenv("OUTPUT_BUCKET"), err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", json_filename, os.Getenv("OUTPUT_BUCKET"))

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	lambda.Start(Handler)
}
