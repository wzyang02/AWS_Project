package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
    "fmt"
    "os"
)

func createS3() {
    if len(os.Args) != 2 {
        exitErrorf("Bucket name missing!\nUsage: %s bucket_name", os.Args[0])
    }

    bucket := os.Args[1]

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    svc := s3.New(sess)

    // Create the S3 Bucket
    _, err = svc.CreateBucket(&s3.CreateBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Unable to create bucket %q, %v", bucket, err)
    }
    fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

    err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
    }

    fmt.Printf("Bucket %q successfully created\n", bucket)
}

func uploadObject() {
    if len(os.Args) != 3 {
        exitErrorf("bucket and file name required\nUsage: %s bucket_name filename",
            os.Args[0])
    }

    bucket := os.Args[1]
    filename := os.Args[2]

    file, err := os.Open(filename)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", filename, err)
    }

    defer file.Close()

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    // Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
    // for more information on configuring part size, and concurrency.
    //
    // http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
    uploader := s3manager.NewUploader(sess)

    // Upload the file's body to S3 bucket as an object with the key being the
    // same as the filename.
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),

        // Can also use the `filepath` standard library package to modify the
        // filename as need for an S3 object key. Such as turning absolute path
        // to a relative path.
        Key: aws.String(filename),

        // The file to be uploaded. io.ReadSeeker is preferred as the Uploader
        // will be able to optimize memory when uploading large content. io.Reader
        // is supported, but will require buffering of the reader's bytes for
        // each part.
        Body: file,
    })
    if err != nil {
        // Print the error and exit.
        exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
    }

    fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

func downloadObject() {
    if len(os.Args) != 3 {
        exitErrorf("Bucket and item names required\nUsage: %s bucket_name item_name",
            os.Args[0])
    }

    bucket := os.Args[1]
    item := os.Args[2]

    file, err := os.Create(item)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", item, err)
    }

    defer file.Close()

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    downloader := s3manager.NewDownloader(sess)

    numBytes, err := downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(item),
        })
    if err != nil {
        exitErrorf("Unable to download item %q, %v", item, err)
    }

    fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func deleteObject() {
    if len(os.Args) != 3 {
        exitErrorf("Bucket and object name required\nUsage: %s bucket_name object_name",
            os.Args[0])
    }

    bucket := os.Args[1]
    obj := os.Args[2]

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    // Create S3 service client
    svc := s3.New(sess)

    // Delete the item
    _, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(obj)})
    if err != nil {
        exitErrorf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
    }

    err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(obj),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for object %q to be deleted, %v", obj, err)
    }

    fmt.Printf("Object %q successfully deleted\n", obj)
}

func deleteBucket() {
    if len(os.Args) != 2 {
        exitErrorf("bucket name required\nUsage: %s bucket_name", os.Args[0])
    }

    bucket := os.Args[1]

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    // Create S3 service client
    svc := s3.New(sess)

    // Delete the S3 Bucket
    // It must be empty or else the call fails
    _, err = svc.DeleteBucket(&s3.DeleteBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Unable to delete bucket %q, %v", bucket, err)
    }

    // Wait until bucket is deleted before finishing
    fmt.Printf("Waiting for bucket %q to be deleted...\n", bucket)

    err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        exitErrorf("Error occurred while waiting for bucket to be deleted, %v", bucket)
    }

    fmt.Printf("Bucket %q successfully deleted\n", bucket)
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}