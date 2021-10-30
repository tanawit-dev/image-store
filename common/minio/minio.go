package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tanawit-dev/image-store/config"
	"io/ioutil"
	"log"
)

const BucketName = "image-storage"
const Location = "ap-southeast-1"
const contentType = "image/jpeg"

var minioClient *minio.Client
var ctx = context.Background()

func Init() {
	minioConfig := config.GetConfig().Minio
	endpoint := fmt.Sprintf("%s:%s", minioConfig.Host, minioConfig.Port)

	mc, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	minioClient = mc

	log.Printf("mino has started suceesfully")

	createBucket()
}

func StoreImage(fileName string, content []byte) error {
	reader := bytes.NewReader(content)
	_, err := minioClient.PutObject(
		ctx,
		BucketName, fileName, reader, reader.Size(), minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Println("upload image successfully")

	return nil
}

func GetImage(fileName string) ([]byte, error) {
	object, err := minioClient.GetObject(ctx, BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	imageContent, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return imageContent, nil
}

func createBucket() {
	err := minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: Location})
	if err != nil {
		exists, err := minioClient.BucketExists(ctx, BucketName)
		if err == nil && exists {
			log.Printf("bucket %s already existed", BucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s \n", BucketName)
	}
}
