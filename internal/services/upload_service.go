package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/chukwiz/go-shop/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type %s", ext)
	}

	newFileName := uuid.New().String()

	path := fmt.Sprintf("products/%d/%s%s", productID, newFileName, ext)

	return s.provider.UploadFile(file, path)
}

func isValidImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}
