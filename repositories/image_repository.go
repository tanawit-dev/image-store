package repositories

import (
	"github.com/tanawit-dev/image-store/common/db"
	"github.com/tanawit-dev/image-store/models"
	"log"
)

func CreateImage(image *models.Image) error {
	tx := db.GetDB().Create(image)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
		return tx.Error
	}

	return nil
}

func FindImageById(id uint) (models.Image, error) {
	var imageModel models.Image
	tx := db.GetDB().First(&imageModel, id)

	return imageModel, tx.Error
}
