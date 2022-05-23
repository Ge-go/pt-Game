package datamodels

import "gorm.io/gorm"

type ConfirmUpload struct {
	Id              int64          `gorm:"column:id;type:bigint(20);not null;primary_key;AUTO_INCREMENT"`
	MateId          int64          `gorm:"column:mate_id;not null;type:bigint(20);index"`
	MateFileName    string         `gorm:"column:mate_file_name;not null;type:varchar(255)"`
	MateFileFormat  string         `gorm:"column:mate_file_format;not null;type:varchar(20)"`
	ThumFileFormat  string         `gorm:"column:thum_file_format;not null;type:varchar(20)"`
	MateFileSize    string         `gorm:"column:mate_file_size;not null;type:varchar(20)"`
	MateType        string         `gorm:"column:mate_type;not null;type:varchar(255)"`
	MateModule      string         `gorm:"column:mate_module;not null;type:varchar(255)"`
	MatePermission  int            `gorm:"column:mate_permission;not null;type:bigint(20);index"`
	MateVersion     int64          `gorm:"column:mate_version;not null;type:bigint(20);index"`
	MatePath        string         `gorm:"column:mate_path;not null;type:varchar(255);index"`
	ThumPath        string         `gorm:"column:thum_path;not null;type:varchar(255);index"`
	MateUuid        string         `gorm:"column:mate_uuid;not null;type:varchar(255);index"`
	Language        string         `gorm:"column:language;not null;type:varchar(255);index"`
	Uploader        string         `gorm:"column:uploader;not null;type:varchar(255)"`
	MateStatus      int            `gorm:"column:mate_status;not null;type:tinyint(1);default 0;index"`
	CronReleaseTime int64          `gorm:"column:cron_release_time;not null;type:bigint(20)"`
	CreateTime      int64          `gorm:"column:create_time;not null;type:bigint(20)"`
	UpdateTime      int64          `gorm:"column:update_time;not null;type:bigint(20)"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;not null;type:timestamp"`
}

func (ConfirmUpload) TableName() string {
	return "admin_confirm_upload"
}

type MateDownload struct {
	Id          int64  `gorm:"column:id;type:bigint(20);not null;primary_key;AUTO_INCREMENT"`
	UserId      int64  `gorm:"column:user_id;not null;type:bigint(20);index"`
	MateId      int64  `gorm:"column:mate_id;not null;type:bigint(20);index"`
	MatePath    string `gorm:"column:mate_path;not null;type:varchar(255)"`
	MateUuid    string `gorm:"column:mate_uuid;not null;type:varchar(255);index"`
	MateVersion int64  `gorm:"column:mate_version;not null;type:bigint(20)"`
	CreateTime  int64  `gorm:"column:create_time;not null;type:bigint(20)"`
}

func (MateDownload) TableName() string {
	return "streamer_mate_download"
}

type MateFavorite struct {
	Id         int64          `gorm:"column:id;type:bigint(20);not null;primary_key;AUTO_INCREMENT"`
	UserId     int64          `gorm:"column:user_id;not null;type:bigint(20);index"`
	MateId     int64          `gorm:"column:mate_id;not null;type:bigint(20);index"`
	MateUuid   string         `gorm:"column:mate_uuid;not null;type:varchar(255);index"`
	CreateTime int64          `gorm:"column:create_time;not null;type:bigint(20)"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;not null;type:timestamp"`
}

func (MateFavorite) TableName() string {
	return "streamer_mate_favorite"
}
