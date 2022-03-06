package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanawit-dev/image-store/services"
	"github.com/tanawit-dev/image-store/utils"
	"mime/multipart"
	"net/http"
	"strconv"
)

const ImageContentType = "image/jpeg"

type ImageController struct {
	Service services.ImageService
}

func ProvideImageController(service services.ImageService) ImageController {
	return ImageController{Service: service}
}

func (con ImageController) GetImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	imageContent, err := con.Service.GetImage(uint(id))
	if err != nil {
		if utils.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("not found image id: %v", id)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.Data(http.StatusOK, ImageContentType, imageContent)
}

func (con ImageController) UploadImage(c *gin.Context) {
	request := UploadImageRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdImage, err := con.Service.SaveImage(request.Image, request.Uploader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK, UploadImageResponse{
			Url:      createdImage.GetUrl(),
			Uploader: createdImage.Uploader,
		},
	)
}

type UploadImageRequest struct {
	Image    *multipart.FileHeader `form:"file" binding:"required"`
	Uploader string                `form:"uploader" binding:"required"`
}

type UploadImageResponse struct {
	Url      string `json:"url"`
	Uploader string `json:"uploader"`
}
