package mysql

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/config"
)

// New new db
func New() *gorm.DB {
	db, err := gorm.Open("mysql", viper.GetString(config.KeyMysqlDSN))
	if err != nil {
		panic(errors.WithMessage(err, "mysql.New().Open() error"))
	}
	db.SingularTable(true)
	db.AutoMigrate(&model.CheckinRecord{}, &model.Customer{}, &model.IssueRecord{}, &model.Merchant{}, &model.User{})
	return db
}

// IsError 是否是真的错误，并不是未找到数据
func IsError(err error) bool {
	return err != nil && err != gorm.ErrRecordNotFound && err != sql.ErrNoRows
}
