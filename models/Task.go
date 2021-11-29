package models

type Task struct {
	ID         int64  `json:"id"`
	FolderId   string `json:"folderId"`
	FolderName string `json:"folderName"`
	Status     string `json:"status"`
}
