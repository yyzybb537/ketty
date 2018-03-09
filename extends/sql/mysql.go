package ketty

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yyzybb537/ketty/extends/log"
	"sync/atomic"
	"time"
)

const LogMysql = "LogMysql"

type MySQLConfig struct {
	Address string
	Dbname  string
	User    string
	Pwd     string
}

type sqlDB struct {
	*sql.DB
}

type Rows struct {
	*sql.Rows
}

type Result struct {
	sql.Result
}

type DB struct {
	*sqlDB
	Address    string
	using      int64
	DbNameBase string
}

func NewDB(ori *sql.DB, addr string) *DB {
	return &DB{
		sqlDB:   &sqlDB{ori},
		Address: addr,
	}
}

func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	nowUsing := atomic.AddInt64(&db.using, 1) - 1
	log.S(LogMysql).Infof("Begin(Using=%d) MySQL.Query @%s TraceID:%s run: %s",
		nowUsing, db.Address, query)
	start := time.Now()
	rows, err := db.sqlDB.Query(query, args...)
	costT := time.Since(start)
	cost := costT.String()
	nowUsing = atomic.AddInt64(&db.using, -1)
	if err != nil {
		log.S(LogMysql).Errorf("End(Using=%d) MySQL.Query @%s TraceID:%s run: %s. cost:%s. err=%+v",
			nowUsing, db.Address, query, cost, err)
	} else {
		log.S(LogMysql).Infof("End(Using=%d) MySQL.Query @%s TraceID:%s run: %s. cost:%s.",
			nowUsing, db.Address, query, cost)
	}
	return &Rows{rows}, err
}

func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	nowUsing := atomic.AddInt64(&db.using, 1) - 1
	log.S(LogMysql).Infof("Begin(Using=%d) MySQL.Exec @%s TraceID:%s run: %s",
		nowUsing, db.Address, query)
	start := time.Now()
	result, err := db.sqlDB.Exec(query, args...)
	costT := time.Since(start)
	cost := costT.String()
	nowUsing = atomic.AddInt64(&db.using, -1)
	if err != nil {
		log.S(LogMysql).Errorf("End(sing=%d) MySQL.Exec @%s TraceID:%s run: %s. cost:%s. err=%+v",
			nowUsing, db.Address, query, cost, err)
	} else {
		log.S(LogMysql).Infof("End(Using=%d) MySQL.Exec @%s TraceID:%s run: %s. cost:%s.",
			nowUsing, db.Address, query, cost)
	}
	return Result{result}, err
}

func MySQLConnect(address, user, pwd, db_name string) (db *DB, err error) {
	return MySQLConnectWithConnNum(address, user, pwd, db_name, 128)
}

func MySQLConnectS(conf MySQLConfig) (db *DB, err error) {
	return MySQLConnect(conf.Address, conf.User, conf.Pwd, conf.Dbname)
}

func MySQLConnectWithConnNum(address, user, pwd, db_name string, num_conn int) (db *DB, err error) {
	dsn := user + ":" + pwd + "@tcp(" + address + ")/" + db_name + "?parseTime=true&interpolateParams=true"
	c, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	c.SetConnMaxLifetime(300 * time.Second)
	c.SetMaxIdleConns(num_conn)
	c.SetMaxOpenConns(num_conn)
	return NewDB(c, address), nil
}
