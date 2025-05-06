package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"lazydebrid/internal/models"
)

const (
	BaseURL                = "https://api.real-debrid.com/rest/1.0"
	TorrentsEndpointURL    = BaseURL + "/torrents"
	DownloadsURL           = BaseURL + "/downloads?page=1&limit=4990"
	TorrentsURL            = BaseURL + "/torrents/"
	TorrentsAddMagnetURL   = TorrentsEndpointURL + "/addMagnet/"
	TorrentsStatusURL      = TorrentsEndpointURL + "/info/"
	TorrentsDeleteURL      = TorrentsEndpointURL + "/delete/"
	TorrentsSelectFilesURL = TorrentsEndpointURL + "/selectFiles/"
)

func UnrestrictLink(link string) (models.TorrentFileDetailed, error) {
	var file models.TorrentFileDetailed
	var apiErr struct {
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
	}

	data := url.Values{"link": {link}}
	body, err := PostForm(BaseURL+"/unrestrict/link/", data)
	if err != nil {
		return file, fmt.Errorf("request failed: %w", err)
	}

	// Try parsing as success type
	if err := json.Unmarshal(body, &file); err == nil && file.Filename != "" {
		return file, nil
	}

	// Fallback to parse API error
	if err := json.Unmarshal(body, &apiErr); err == nil {
		return file, fmt.Errorf("API error: %s (code %d)", apiErr.Error, apiErr.ErrorCode)
	}

	return file, fmt.Errorf("unrecognized response: %s", string(body))
}
