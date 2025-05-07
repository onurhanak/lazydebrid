package actions

import (
	"lazydebrid/internal/models"
)

func RemoveID(slice []models.ActiveDownload, item string) []models.ActiveDownload {
	result := []models.ActiveDownload{}
	for _, v := range slice {
		if v.ID != item {
			result = append(result, v)
		}
	}
	return result
}
