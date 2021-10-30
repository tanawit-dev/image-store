package services

import (
	"bytes"
	"fmt"
	"github.com/tanawit-dev/image-store/common/db"
	"github.com/tanawit-dev/image-store/common/minio"
	"github.com/tanawit-dev/image-store/models"
	"github.com/tanawit-dev/image-store/repositories"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
)

const FileExtension = "jpeg"

func GetImage(id uint) ([]byte, error) {
	imageModel, err := repositories.FindImageById(id)
	if err != nil {
		return nil, err
	}

	imageContent, err := minio.GetImage(generateFileName(imageModel))
	if err != nil {
		return nil, err
	}

	return imageContent, nil
}

func SaveImage(image *multipart.FileHeader, uploader string) (*models.Image, error) {
	newImage := models.Image{
		FileName: image.Filename,
		Uploader: uploader,
	}

	err := db.GetDB().Transaction(
		func(tx *gorm.DB) error {

			if err := repositories.CreateImage(&newImage); err != nil {
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

			fileName := generateFileName(newImage)
			if err := minio.StoreImage(fileName, buf.Bytes()); err != nil {
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
