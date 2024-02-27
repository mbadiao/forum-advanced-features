package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func GetCategory(CheckboxValues []string) []int {
	var CategoriesId []int
	for _, v := range CheckboxValues {
		if v == "All" {
			CategoriesId = []int{1, 2, 3, 4, 5}
		}
		if v == "Tech" {
			CategoriesId = append(CategoriesId, 1)
		}
		if v == "Actu" {
			CategoriesId = append(CategoriesId, 2)
		}
		if v == "Mode" {
			CategoriesId = append(CategoriesId, 3)
		}
		if v == "Sport" {
			CategoriesId = append(CategoriesId, 4)
		}
		if v == "Edu" {
			CategoriesId = append(CategoriesId, 5)
		}
	}
	return CategoriesId
}

func Checkcategory(category []string) bool {
	Mapcategory := map[string]bool{
		"All":   true,
		"Tech":  true,
		"Actu":  true,
		"Mode":  true,
		"Sport": true,
		"Edu":   true,
	}
	found := true
	for _, v := range category {
		if !CompareCategory(Mapcategory, v) {
			found = false
		}
	}
	return found
}

func CompareCategory(categoriesMap map[string]bool, categoryToCheck string) bool {
	_, found := categoriesMap[categoryToCheck]
	return found
}

func IsValidImage(file multipart.File, handler *multipart.FileHeader) bool {
	//contentType := handler.Header.Get("Content-Type")
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	// Reset file offset to beginning for subsequent reads
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return false
	}

	// Detect content type based on the first 512 bytes
	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/svg+xml", "image/jpeg", "image/gif", "image/png":
		// OK, continuez le traitement
	default:
		fmt.Println(contentType)
		fmt.Println("type")
		return false
	}
	return true
}
