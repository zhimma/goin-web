package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhimma/goin-web/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDbR() *gorm.DB
	GetDbW() *gorm.DB
	DbRClose() error
	DbWClose() error
}

type dbRepo struct {
	DbR *gorm.DB
	DbW *gorm.DB
}

func (d dbRepo) i() {}

func (d dbRepo) GetDbR() *gorm.DB {
	return d.DbR
}

func (d dbRepo) GetDbW() *gorm.DB {
	return d.DbW
}

func (d dbRepo) DbRClose() error {
	sqlDB, err := d.DbR.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d dbRepo) DbWClose() error {
	sqlDB, err := d.DbW.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// 1.入口方法，返回repo，repo中包含2个数据库资源，读和写 及2个关闭资源的方法
func New() (Repo, error) {
	cfg := config.Get().Mysql

	dbr, err := dbConnect(cfg.Read.Username, cfg.Read.Password, cfg.Read.Host, cfg.Read.Database)
	if err != nil {
		return nil, err
	}

	dbw, err := dbConnect(cfg.Write.Username, cfg.Write.Password, cfg.Write.Host, cfg.Write.Database)
	if err != nil {
		return nil, err
	}

	// 获取读库的链接实例
	return &dbRepo{
		DbR: dbr,
		DbW: dbw,
	}, nil
}

func dbConnect(user, pass, addr, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info), // 日志配置
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", dbName))
	}
	// 设置字符集
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	cfg := config.Get().Mysql.Base

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConn)

	// 设置最大连接超时
	sqlDb.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifeTime)

	// 使用插件

	return db, nil
}
