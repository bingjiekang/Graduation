package common

import "time"

type FileUploadAndDownload struct {
	ID        int       `gorm:"primarykey"`                                          // 主键ID
	UUid      int64     `json:"uUid" form:"uUid" gorm:"column:u_uid;comment:唯一标识id"` // 管理员唯一标识id
	Name      string    `json:"name" gorm:"comment:文件名"`                             // 文件名
	Url       string    `json:"url" gorm:"comment:文件地址"`                             // 文件路径
	Tag       string    `json:"tag" gorm:"comment:文件标签"`                             // 文件标签(场景)
	Key       string    `json:"key" gorm:"comment:编号"`                               // 编号
	CreatedAt time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallShoppingCartItem 表名
func (FileUploadAndDownload) TableName() string {
	return "file_upload_and_download"
}

// file struct, 文件结构体
type ExaFile struct {
	ID           int
	FileName     string
	FileMd5      string
	FilePath     string
	ExaFileChunk []ExaFileChunk
	ChunkTotal   int
	IsFinish     bool
	CreateAt     time.Time
	UpdateAt     time.Time
}

// file chunk struct, 切片结构体
type ExaFileChunk struct {
	ID              int
	ExaFileID       uint
	FileChunkNumber int
	FileChunkPath   string
	CreateAt        time.Time
	UpdateAt        time.Time
}
