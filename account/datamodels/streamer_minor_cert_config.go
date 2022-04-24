package datamodels

import "time"

type MinorCertConfig struct {
	Id             int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT"`                   // ID
	ShortName      string    `gorm:"column:short_name;type:varchar(150);NOT NULL"`                        // 国家英语短名称
	AlphaTwoCode   string    `gorm:"column:alpha_two_code;type:varchar(10);NOT NULL"`                     // alphaCODE
	AlphaThreeCode string    `gorm:"column:alpha_three_code;type:varchar(15);NOT NULL"`                   // alphaCode
	Numeric        string    `gorm:"column:numeric;type:varchar(15);NOT NULL"`                            // 数字编号
	CertType       int       `gorm:"column:cert_type;type:tinyint(1)"`                                    // 认证类型
	AdultAge       int       `gorm:"column:adult_age;type:tinyint(2)"`                                    // 成年标准依据
	IsEuropean     int       `gorm:"column:is_european;type:tinyint(1)"`                                  // 是否欧盟国家 0 不是 1 是 2不是但是按欧盟标准
	Creater        int       `gorm:"column:creater;type:int(11);NOT NULL"`                                // 创建人ID
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
}

func (MinorCertConfig) TableName() string {
	return "streamer_minor_cert_config"
}
