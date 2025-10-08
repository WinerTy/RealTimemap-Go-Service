package entity

import (
	"fmt"
)

// Image структура для представления JSON поля картинки из БД
type Image struct {
	ContentPath   string   `json:"content_path"`
	FileName      string   `json:"filename"`
	ContentType   string   `json:"content_type"`
	Size          int      `json:"size"`
	Files         []string `json:"files"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	FileID        string   `json:"file_id"`
	UploadStorage string   `json:"upload_storage"`
}

//[{"content_path": null, "filename":
//"5. \u0411\u0435\u043b\u044b\u0439 \u0442\u0438\u0433\u0440\u0435\u043d\u043e\u043a.jpg",
//"content_type": "image/jpeg", "size": 1126898, "files": ["marks/5a7760bb-f589-4f54-92fd-f4354f041fe3"],
//"width": 2560, "height": 1600, "file_id": "5a7760bb-f589-4f54-92fd-f4354f041fe3", "upload_storage":
//"marks", "uploaded_at": "2025-09-29T13:03:15.909070", "path": "marks/5a7760bb-f589-4f54-92fd-f4354f041fe3",
//"url": "uploads\\marks\\5a7760bb-f589-4f54-92fd-f4354f041fe3", "saved": true}]

func (i *Image) BuildUrl() string {
	if i.FileID != "" && i.UploadStorage != "" {
		return fmt.Sprintf("https://realtimemap.ru/media/%s/%s", i.UploadStorage, i.FileID)
	}
	return ""
}
