package repositories

import (
	"github.com/tanawit-dev/image-store/models"
	"gorm.io/gorm"
	"log"
)

type ImageRepository struct {
	db *gorm.DB
}

func ProvideImageRepository(db *gorm.DB) ImageRepository {
	return ImageRepository{db: db}
}

func (repo ImageRepository) Create(image *models.Image) error {
	tx := repo.db.Create(image)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
		return tx.Error
	}

	return nil
}

func (repo ImageRepository) FindById(id uint) (models.Image, error) {
	var imageModel models.Image
	tx := repo.db.First(&imageModel, id)

	return imageModel, tx.Error
}
