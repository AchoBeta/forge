package database

import (
	"forge/infra/configs"
	"forge/pkg/log/zlog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDataBases 初始化
func initMysql(config configs.IConfig) error {
	dsn := config.GetDBConfig().Dsn
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		zlog.Panicf("MySQL无法连接数据库！: %v", err)
		return err
	}
	zlog.Infof("MySQL连接数据库成功！")
	db = _db

	return nil
}
