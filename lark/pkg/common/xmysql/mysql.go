package xmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"lark/pkg/common/xlog"
	"lark/pkg/conf"
	"log"
	"os"
	"time"
)

var (
	ERR_DB_INSTANCE_IS_EMPTY = errors.New("database instance is empty")
)

var (
	cli *MysqlClient
)

type MysqlClient struct {
	db        *gorm.DB
	cfg       *conf.Mysql
	connected bool
}

func NewMysqlClient(cfg *conf.Mysql) *MysqlClient {
	var (
		err error
	)
	cli = &MysqlClient{cfg: cfg}
	cli.db, err = ConnectDB(cfg)
	cli.connected = err == nil
	return cli
}

func GetDB() *gorm.DB {
	if cli.db == nil {
		var (
			err error
		)
		cli.db, err = ConnectDB(cli.cfg)
		cli.connected = err == nil
	}
	return cli.db
}

func GetTX() *gorm.DB {
	return GetDB().Begin()
}

// 事务处理
func Transaction(handle func(tx *gorm.DB) (err error)) (err error) {
	var (
		db   *gorm.DB
		terr error
	)
	db = GetDB()
	if db == nil {
		err = ERR_DB_INSTANCE_IS_EMPTY
		return
	}
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	err = handle(tx)
	if err != nil {
		terr = tx.Rollback().Error
		if terr != nil {
			xlog.Error(terr.Error())
		}
		return
	}
	err = tx.Commit().Error
	return
}

func ConnectDB(cfg *conf.Mysql) (db *gorm.DB, err error) {
	var (
		dsn   string
		opts  *gorm.Config
		sqlDB *sql.DB
	)
	dsn = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.Db)

	// 定义日志配置
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                   // 慢查询阈值
			LogLevel:                  logger.LogLevel(cfg.LogLevel), // 日志级别
			IgnoreRecordNotFoundError: true,                          // 忽略记录未找到错误
			Colorful:                  false,                         // 禁用彩色打印
		},
	)
	opts = &gorm.Config{
		SkipDefaultTransaction: false, // 禁用默认事务(true: Error 1295: This command is not supported in the prepared statement protocol yet)
		PrepareStmt:            false, // 创建并缓存预编译语句(true: Error 1295)
		Logger:                 gormLogger,
	}

	db, err = gorm.Open(mysql.Open(dsn), opts)
	if err != nil {
		xlog.Error(err.Error())
		return
	}

	sqlDB, err = db.DB()
	if err != nil {
		xlog.Error(err.Error())
		return
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Millisecond)
	return
}

func GetDbName(db *gorm.DB) string {
	var dbName string
	config, err := db.DB()
	if err != nil {
		fmt.Println("Failed to get database configuration:", err)
		return ""
	}
	err = config.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		fmt.Println("Failed to get database name:", err)
		return ""
	}
	return dbName
}

func Connected() bool {
	if cli == nil {
		return false
	}
	return cli.connected
}
