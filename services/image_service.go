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
	repo         repositories.ImageRepository
	imageStorage minio.ImageStorage
	db           *gorm.DB
}

func ProvideImageService(repo repositories.ImageRepository, store minio.ImageStorage, db *gorm.DB) ImageService {
	return ImageService{repo: repo, imageStorage: store, db: db}
}

func (service ImageService) GetImage(id uint) ([]byte, error) {
	imageModel, err := service.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	imageContent, err := service.imageStorage.Load(generateFileName(imageModel))
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

	err := service.db.Transaction(
		func(tx *gorm.DB) error {
			if err := service.repo.Create(&newImage); err != nil {
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

			if err := service.imageStorage.Store(generateFileName(newImage), buf.Bytes()); err != nil {
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
