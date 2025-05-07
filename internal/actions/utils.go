package actions

import (
	"lazydebrid/internal/models"
	"log"
)

func RemoveItem(slice []models.ActiveDownload, item string) []models.ActiveDownload {
	result := []models.ActiveDownload{}
	for _, v := range slice {
		log.Println(v.ID, item)
		if v.ID != item {
			result = append(result, v)
		}
	}
	return result
}
