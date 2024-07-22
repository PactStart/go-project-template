package relation

import (
	"fmt"
	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"time"
)

const (
	maxRetry = 100 // number of retries
	Driver   = "mysql"
)

type option struct {
	Username      string
	Password      string
	Address       []string
	Database      string
	LogLevel      int
	SlowThreshold int
	MaxLifeTime   int
	MaxOpenConn   int
	MaxIdleConn   int
	Connect       func(dsn string, maxRetry int) (*gorm.DB, error)
}

// newMysqlGormDB Initialize the database connection.
func newMysqlGormDB(o *option) (*gorm.DB, error) {
	err := maybeCreateDatabase(o)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		o.Username, o.Password, o.Address[0], o.Database)
	sqlLogger := log.NewSqlLogger(
		logger.LogLevel(o.LogLevel),
		true,
		time.Duration(o.SlowThreshold)*time.Millisecond,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(o.MaxLifeTime))
	sqlDB.SetMaxOpenConns(o.MaxOpenConn)
	sqlDB.SetMaxIdleConns(o.MaxIdleConn)
	return db, nil
}

// maybeCreateDatabase creates a database if it does not exists.
func maybeCreateDatabase(o *option) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		o.Username, o.Password, o.Address[0], "mysql")

	var db *gorm.DB
	var err error
	if f := o.Connect; f != nil {
		db, err = f(dsn, maxRetry)
	} else {
		db, err = connectToDatabase(dsn, maxRetry)
	}
	if err != nil {
		panic(err.Error() + " Open failed " + dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	sql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` default charset utf8mb4 COLLATE utf8mb4_unicode_ci",
		o.Database,
	)
	err = db.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("init db %w", err)
	}
	return nil
}

// connectToDatabase Connection retry for mysql.
func connectToDatabase(dsn string, maxRetry int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i <= maxRetry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), nil)
		if err == nil {
			return db, nil
		}
		if mysqlErr, ok := err.(*mysqldriver.MySQLError); ok && mysqlErr.Number == 1045 {
			return nil, err
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	return nil, err
}

// NewGormDB gorm mysql.
func NewGormDB() (*gorm.DB, error) {
	errs.AddReplace(gorm.ErrRecordNotFound, errs.ErrRecordNotFound)
	errs.AddErrHandler(replaceDuplicateKey)

	return newMysqlGormDB(&option{
		Username:      config.Config.Mysql.Username,
		Password:      config.Config.Mysql.Password,
		Address:       config.Config.Mysql.Address,
		Database:      config.Config.Mysql.Database,
		LogLevel:      config.Config.Mysql.LogLevel,
		SlowThreshold: config.Config.Mysql.SlowThreshold,
		MaxLifeTime:   config.Config.Mysql.MaxLifeTime,
		MaxOpenConn:   config.Config.Mysql.MaxOpenConn,
		MaxIdleConn:   config.Config.Mysql.MaxIdleConn,
	})
}

func replaceDuplicateKey(err error) errs.CodeError {
	if IsMysqlDuplicateKey(err) {
		return errs.ErrDuplicateKey
	}
	return nil
}

func IsMysqlDuplicateKey(err error) bool {
	if mysqlErr, ok := err.(*mysqldriver.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
