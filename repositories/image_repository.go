package repositories

import (
	"github.com/tanawit-dev/image-store/common/db"
	"github.com/tanawit-dev/image-store/models"
	"log"
)

func CreateImage(image *models.Image) error {
	create := db.GetDB().Create(image)
	if create.Error != nil {
		log.Fatalln(create.Error)
		return create.Error
	}

	return nil
}

func FindImageById(id uint) (models.Image, error) {
	var imageModel models.Image
	tx := db.GetDB().First(&imageModel, id)

	return imageModel, tx.Error
}
