package models

type DebridDownload struct {
	Id         string `json:"id"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mimeType"`
	Filesize   int64  `json:"filesize"`
	Link       string `json:"link"`
	Host       string `json:"host"`
	HostIcon   string `json:"host_icon"`
	Chunks     int64  `json:"chunks"`
	Download   string `json:"download"`
	Streamable int64  `json:"streamable"`
	Generated  string `json:"generated"`
}

type TorrentStatus struct {
	ID               string        `json:"id"`
	Filename         string        `json:"filename"`
	OriginalFilename string        `json:"original_filename"`
	Hash             string        `json:"hash"`
	Bytes            int64         `json:"bytes"`
	OriginalBytes    int64         `json:"original_bytes"`
	Host             string        `json:"host"`
	Split            int           `json:"split"`
	Progress         int           `json:"progress"`
	Status           string        `json:"status"`
	Added            string        `json:"added"`
	Files            []TorrentFile `json:"files"`
	Links            []string      `json:"links"`
}

type TorrentFile struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Bytes    int64  `json:"bytes"`
	Selected int    `json:"selected"`
}

type ActiveDownload struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}
