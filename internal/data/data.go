package data

import "lazydebrid/internal/models"

var (
	UserDownloads   []models.Torrent
	DownloadMap     = make(map[string]models.Torrent)
	ActiveDownloads []models.ActiveDownload
	FilesMap        = make(map[string]models.Download)
)

type LineMapping struct {
	ID string
}

var TorrentLineIndex []string
