package config

import "sync"

type db struct {
	mu       sync.Mutex
	Driver   DbDriver
	Host     string
	Username string
	Password string
	Database string
}

type DbDriver int

const (
	MysqlDriver DbDriver = iota + 1
	PostgreSqlDriver
)

func (dbd DbDriver) String() string {
	if dbd == 1 {
		return "mysql"
	}

	if dbd == 2 {
		return "postgresql"
	}

	return ""
}

func (db *db) SetHost(h string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Host = h
}

func (db *db) SetUsername(u string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Username = u
}

func (db *db) SetPassword(p string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Password = p
}

func (db *db) SetDbName(d string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Database = d
}

func (db *db) SetDriver(d DbDriver) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Driver = d
}
