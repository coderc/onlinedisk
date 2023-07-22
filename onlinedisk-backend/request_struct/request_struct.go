package request

type RequestUserInfo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type MultipleInfo struct {
	UploadId   string `json:"uploadId"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	ChunkSize  int64  `json:"chunkSize" binding:"omitempty"`
	ChunkCount int64  `json:"chunkCount"`
	FileSHA1   string `json:"fileSHA1"`
}
