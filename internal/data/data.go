package data

import "lazydebrid/internal/models"

var (
	UserDownloads = make(map[int]models.Torrent)
	// DownloadMap       = make(map[string]models.Torrent)
	ActiveDownloads []models.ActiveDownload
	// ActiveDownloadsMap = make(map[int]models.Torrent)
	FilesMap = make(map[string]models.TorrentFileDetailed)
)
