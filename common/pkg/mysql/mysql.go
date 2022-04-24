package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ptc-Game/common/pkg/config"
	"time"
)

func New(config config.MySQLConfig) (*gorm.DB, error) {
	db, err := newMysql(config)
	return db, err
}

func newMysql(config config.MySQLConfig) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		config.UserName,
		config.Password,
		config.Addr,
		config.Name,
		true,
		"UTC")

	gormConfig := &gorm.Config{}
	// 根据配置文件,是否打印sql日志
	if config.ShowLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dns), gormConfig)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to open mysql")
	}
	sql, err := db.DB()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to config mysql")
	}
	sql.SetMaxIdleConns(config.MaxIdleConn)
	sql.SetMaxOpenConns(config.MaxOpenConn)
	sql.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifeTime))

	return db, err
}
