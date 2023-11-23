package before

import (
	"errors"
	"time"

	dbtype "github.com/tzRex/freely-handle/enum/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// gorm 连接配置
type DatabaseConnetStu struct {
	DbType         string        // 数据库类型：mysql、mssql、oracle
	DbHost         string        // 连接主体：127.0.0.1
	DbPort         string        // 端口号
	DbUsername     string        // 数据库用户账号
	DbPassword     string        // 数据库用户密码
	DbName         string        // 数据库名称
	DbMaxOpenConns int           // 最大连接数
	DbMaxIdleConns int           // 最大空闲连接数量
	DbMaxLifetime  time.Duration // 最大生存时间
}

var Database *gorm.DB

/**
 * 数据库连接
 * @Return 是否连接成功
 */
func ConnectionDatabase(config *DatabaseConnetStu) error {
	switch config.DbType {
	case dbtype.Mysql:
		return connectMysql(config)
	case dbtype.Mssql:
		return connectMssql(config)
	default:
		return errors.New("databse.not.find")
	}
}

/**
 * 连接Mysql数据库
 */
func connectMysql(config *DatabaseConnetStu) error {
	driverDns := config.DbUsername + ":" + config.DbPassword +
		"@tcp(" + config.DbHost + ":" + config.DbPort + ")/" +
		config.DbName +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(
		mysql.Open(driverDns),
		&gorm.Config{
			PrepareStmt: true,
		},
	)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(config.DbMaxIdleConns)

	sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)

	sqlDB.SetConnMaxLifetime(config.DbMaxLifetime)

	return nil
}

/**
 * 连接Mssql数据库
 */
func connectMssql(config *DatabaseConnetStu) error {
	return nil
}

/**
 * 获取数据库
 */
func GetDB(tx ...*gorm.DB) *gorm.DB {
	var db *gorm.DB

	if len(tx) > 0 {
		db = tx[0]
	} else {
		db = Database
	}

	return db.Session(&gorm.Session{
		NewDB:           true,
		PrepareStmt:     true,
		CreateBatchSize: 500,
	})
}

/**
 * 数据库事务
 */
func DBTransaction(call func(tx *gorm.DB) error) error {
	return Database.Transaction(call)
}
