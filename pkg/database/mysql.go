package database

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

func MysqlOpenConnection(user, pass, host_port, db string) (*sqlx.DB, error) {
	conn, errOpen := sqlx.Open(`mysql`, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", user, pass, host_port, db))
	if errOpen != nil {
		return nil, errOpen
	}

	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(30)
	conn.SetConnMaxLifetime(time.Minute * 10)

	if errPing := conn.Ping(); errPing != nil {
		return nil, errPing
	}

	return conn, nil
}
