package data

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Data struct {
	Db *gorm.DB
}

func (this *Data)DB(ctx context.Context) (*gorm.DB) {
	return this.Db.WithContext(ctx)
}

type Options struct {
	Address string
	UserName string
	Password string
	DBName string
	Logger  logger.Writer
	Charset string
}

func New(opt *Options) (*gorm.DB, error) {
	if opt.Charset == "" {
		opt.Charset = "utf8"
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local", opt.UserName, opt.Password,
		opt.Address, opt.DBName, opt.Charset)
	cfg := &gorm.Config{}
	if opt.Logger != nil {
		cfg.Logger = logger.New(
			opt.Logger, // io writer
			logger.Config{
				SlowThreshold:              time.Second,   // Slow SQL threshold
				LogLevel:                   logger.Info,   // Log level
				IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,          // Disable color
			},
		)
	}
	db, err := gorm.Open(mysql.Open(connStr), cfg)
	if err != nil {
		return nil, err
	}
	instance, err := db.DB()
	if err != nil {
		return nil, err
	}
	instance.SetMaxIdleConns(5)
	instance.SetMaxOpenConns(50)
	instance.SetConnMaxLifetime(time.Hour)

	cfg.Logger.Info(context.Background(), "DB connect success.  " + fmt.Sprintf("address: %s, user: %s, db: %s", opt.Address, opt.UserName, opt.DBName))
	return db, nil
}

