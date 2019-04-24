package models

type FileDownloadAddressResponse struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	FileName     string `json:"fileName"`
	DownloadPath string `json:"downloadPath"`
}
