package datamodels

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id             int            `gorm:"column:id;not null;primaryKey;autoIncrement;"`
	UserName       string         `gorm:"column:user_name;index"`
	Password       string         `gorm:"column:password"`
	Kind           int            `gorm:"column:kind"`
	Email          string         `gorm:"column:email;unique"`
	Gender         int            `gorm:"column:gender"`
	Country        string         `gorm:"column:country;index"`
	Region         string         `gorm:"column:region;index"`
	GoogleSub      string         `gorm:"column:google_sub"`
	YoutubeAccount string         `gorm:"column:youtube_account"`
	YoutubeFansNum int            `gorm:"column:youtube_fans_num"`
	YoutubeChannel string         `gorm:"column:youtube_channel"`
	Language       string         `gorm:"column:language;default:en"`
	Level          string         `gorm:"column:level"`
	Birthday       string         `gorm:"column:birthday;"`
	IsAdult        int            `gorm:"column:is_adult"`
	IsSigned       int            `gorm:"column:is_signed"`
	IsLocked       int            `gorm:"column:is_locked"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null;index"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	Updater        int            `gorm:"column:updater"`
}

//签名表
type StreamerSign struct {
	Id           int    `gorm:"column:id;not null;primaryKey;autoIncrement;"`
	UserId       int    `gorm:"column:user_id"`
	EnvelopeId   string `gorm:"column:envelope_id"`
	Status       string `gorm:"column:status"`
	Manner       string `gorm:"column:manner"`
	CreateTime   string `gorm:"column:create_time"`
	CompleteTime string `gorm:"column:complete_time"` //完成时间
	DeclineTime  string `gorm:"column:decline_time"`  //拒绝时间
	DeclineReson string `gorm:"column:decline_reson"` //拒绝理由
	DownloadUrl  string `gorm:"column:download_url"`
}

func (User) TableName() string {
	return "streamer_user"
}
