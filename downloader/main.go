package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "go.uber.org/automaxprocs"
)

//log "github.com/sirupsen/logrus"

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {

	// =====================
	// Get OS parameter
	// =====================
	var bucket string
	var key string
	var region string
	var multipartsize int64
	var concurrency int

	flag.StringVar(&bucket, "bucket", "", "bucket name")
	flag.StringVar(&key, "key", "", "key name")
	flag.StringVar(&region, "region", "us-east-1", "region name")
	flag.Int64Var(&multipartsize, "multipartsize", 100*1024*1024, "multipart size")
	flag.IntVar(&concurrency, "concurrency", 1, "concurrency")

	flag.Parse()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	s3Svc := s3.New(sess)

	downloader := s3manager.NewDownloaderWithClient(s3Svc, func(d *s3manager.Downloader) {
		d.PartSize = multipartsize
		d.Concurrency = concurrency
	})

	file, _ := os.OpenFile("/dev/null", os.O_RDWR, 0644)

	start := time.Now()
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", key, err)
	}

	duration := time.Since(start).Seconds()

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes", "in", duration, "seconds")
	fmt.Printf("Rate MB/s %f\n", float64(numBytes)/float64(duration)/float64(1024*1024))
	fmt.Printf("Rate Gbit/s %f\n", float64(8)*float64(numBytes)/float64(duration)/float64(1024*1024*1024))
}
