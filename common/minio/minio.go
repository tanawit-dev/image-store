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

type ImageStorage struct {
	minioClient *minio.Client
	ctx         context.Context
}

func (service ImageStorage) Load(fileName string) ([]byte, error) {
	object, err := service.minioClient.GetObject(service.ctx, BucketName, fileName, minio.GetObjectOptions{})
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

func (service ImageStorage) Store(fileName string, content []byte) error {
	reader := bytes.NewReader(content)
	_, err := service.minioClient.PutObject(
		service.ctx,
		BucketName, fileName, reader, reader.Size(), minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Println("upload image successfully")

	return nil
}

func ProvideImageStorage(config config.Configuration) (ImageStorage, error) {
	ctx := context.Background()
	endpoint := fmt.Sprintf("%s:%s", config.Minio.Host, config.Minio.Port)

	mc, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.Minio.AccessKey, config.Minio.SecretKey, ""),
			Secure: false,
		},
	)

	if err != nil {
		log.Fatalln(err)
		return ImageStorage{}, err
	}

	log.Printf("mino has started suceesfully")

	defer createBucket(mc, ctx)

	return ImageStorage{minioClient: mc, ctx: ctx}, nil
}

func createBucket(minioClient *minio.Client, ctx context.Context) {
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
