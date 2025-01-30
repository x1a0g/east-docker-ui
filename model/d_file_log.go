package model

import "gorm.io/gorm"

type FileLog struct {
	Id       string `json:"id" gorm:"column:file_id"`
	CreateAt string `json:"createAt" gorm:"column:create_at"`
	FilePath string `json:"filePath" gorm:"column:file_path"`
	FileSize int64  `json:"fileSize" gorm:"column:file_size"`
	FileName string `json:"fileName" gorm:"column:file_name"`
}

func (*FileLog) TableName() string {
	return "d_file_log"
}

func (f *FileLog) GetFileLog(db *gorm.DB, id string) (FileLog, error) {
	var fileLog FileLog
	err := db.Where("file_id = ?", id).First(&fileLog).Error
	return fileLog, err
}

// 保存
func (f *FileLog) SaveFileLog(db *gorm.DB) error {
	return db.Create(&f).Error
}
