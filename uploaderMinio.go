package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploaderMinio struct {
	interval     time.Duration
	bucket       string
	objectPrefix string
	buffer       *Buffer
	minio        *minio.Client
}

func NewUploaderMinio(
	interval time.Duration,
	buffer *Buffer,
	bucket string,
	objectPrefix string,
	endpoint string,
	minioKey string,
	minioSecret string) *UploaderMinio {
	// sess := session.Must(session.NewSession())
	// svc := s3.New(sess)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioKey, minioSecret, ""),
		Secure: false,
	})

	if err != nil {
		log.Printf("Creating minio client failed: %s", err)
		return nil
	}

	return &UploaderMinio{
		interval:     interval,
		bucket:       bucket,
		objectPrefix: objectPrefix,
		buffer:       buffer,
		minio:        minioClient,
	}
}

func (u *UploaderMinio) RunLoop() {
	ticker := time.NewTicker(u.interval)
	for {
		<-ticker.C
		u.Run()
	}
}

func (u *UploaderMinio) Run() {
	path, err := u.buffer.Rotate()
	if err != nil {
		log.Printf("Rotating a file failed: %s", err)
		return
	}

	var compressedPath string
	for {
		compressedPath, err = u.compressFile(path)
		if err == nil {
			break
		}
		log.Printf("Compressing a file failed: %s", err)
		log.Printf("Retrying in 10 sec")
		time.Sleep(time.Second * 10)
	}

	for {
		err = u.uploadFile(compressedPath)
		if err == nil {
			break
		}
		log.Printf("Uploading a file failed: %s", err)
		log.Printf("Retrying in 10 sec")
		time.Sleep(time.Second * 10)
	}

	log.Printf("Uploading succeeded")

	err = u.deleteFile(compressedPath)
	if err != nil {
		log.Printf("Deleting %s failed: %s", compressedPath, err)
	}
}

func (u *UploaderMinio) compressFile(path string) (string, error) {
	log.Printf("Compressing %s", path)
	err := exec.Command("gzip", path).Run()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.gz", path), nil
}

func (u *UploaderMinio) uploadFile(path string) error {
	log.Printf("Uploading %s", path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	objectName := fmt.Sprintf("%s%s.ltsv.gz", u.objectPrefix, time.Now().UTC().Format("2006/01/02/20060102_150405"))
	log.Printf("PutObject %s", objectName)
	// _, err = u.s3.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String(u.bucket),
	// 	Key:    aws.String(key),
	// 	Body:   f,
	// })

	ctx := context.Background()

	// _, err := u.minio.FPutObject(ctx,
	// 	u.bucket,
	// 	objectName,
	// 	filePath,
	// 	minio.PutObjectOptions{ContentType: contentType})

	_, err_minio :=
		u.minio.PutObject(ctx,
			u.bucket,
			objectName,
			f,                        // file stream object
			-1,                       // object size
			minio.PutObjectOptions{}) //,

	if err_minio != nil {
		log.Fatalln(err)
	}

	if err_minio != nil {
		return err_minio
	}

	return nil
}

func (u *UploaderMinio) deleteFile(path string) error {
	log.Printf("Deleting %s", path)
	return os.Remove(path)
}
