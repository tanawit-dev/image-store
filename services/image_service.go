package services

import (
	"bytes"
	"fmt"
	"github.com/tanawit-dev/image-store/common/minio"
	"github.com/tanawit-dev/image-store/models"
	"github.com/tanawit-dev/image-store/repositories"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
)

const FileExtension = "jpeg"

type ImageService struct {
	Repo         repositories.ImageRepository
	ImageStorage minio.ImageStorage
	DB           *gorm.DB
}

func (service ImageService) GetImage(id uint) ([]byte, error) {
	imageModel, err := service.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	imageContent, err := service.ImageStorage.Load(generateFileName(imageModel))
	if err != nil {
		return nil, err
	}

	return imageContent, nil
}

func (service ImageService) SaveImage(image *multipart.FileHeader, uploader string) (*models.Image, error) {
	newImage := models.Image{
		FileName: image.Filename,
		Uploader: uploader,
	}

	err := service.DB.Transaction(
		func(tx *gorm.DB) error {
			if err := service.Repo.Create(&newImage); err != nil {
				return err
			}

			file, err := image.Open()
			if err != nil {
				return err
			}
			defer func(file multipart.File) {
				if err := file.Close(); err != nil {
					log.Fatalln(err)
				}
			}(file)

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				return err
			}

			if err := service.ImageStorage.Store(generateFileName(newImage), buf.Bytes()); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &newImage, nil
}

func generateFileName(image models.Image) string {
	return fmt.Sprintf("%d.%s", image.ID, FileExtension)
}
