package models

import (
	"fmt"
	"gorm.io/gorm"
	"net/url"
)

type Image struct {
	gorm.Model
	FileName string `gorm:"not_null"`
	Uploader string `gorm:"not_null"`
}

func (i *Image) GetUrl() string {
	imageUrl := url.URL{
		Path:   fmt.Sprintf("/api/images/%d", i.ID),
	}

	return imageUrl.String()
}
