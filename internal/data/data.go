package data

import "lazydebrid/internal/models"

var (
	UserDownloads     []models.Torrent
	DownloadMap       = make(map[string]models.Torrent)
	ActiveDownloads   []models.ActiveDownload
	ActiveDownloadMap = make(map[string]models.ActiveDownload)
	FilesMap          = make(map[string]models.Download)
	TorrentLineIndex  []string
)
