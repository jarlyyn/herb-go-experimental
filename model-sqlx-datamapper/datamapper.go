package datamapper

import (
	"github.com/jmoiron/sqlx"
)

type DB interface {
	SetDB(db *sqlx.DB)
	DB() *sqlx.DB
	TableName(string) string
}

type DataMapper interface {
	DB() *sqlx.DB
	DBTableName() string
}

func NewDB(db *sqlx.DB, prefix string) *PrefixDB {
	return &PrefixDB{
		db:     db,
		prefix: prefix,
	}
}

type PrefixDB struct {
	db     *sqlx.DB
	prefix string
}

func (d *PrefixDB) SetDB(db *sqlx.DB) {
	d.db = db
}

func (d *PrefixDB) DB() *sqlx.DB {
	return d.db
}

func (d *PrefixDB) SetPrefix(prefix string) {
	d.prefix = prefix
}

func (d *PrefixDB) Prefix() string {
	return d.prefix
}

func (d *PrefixDB) TableName(tableName string) string {
	return d.prefix + tableName
}
func New(db DB, table string) *TableDataMapper {
	return &TableDataMapper{
		db:    db,
		table: table,
	}
}

type TableDataMapper struct {
	db    DB
	table string
}

func (t *TableDataMapper) DB() *sqlx.DB {
	return t.db.DB()
}

func (t *TableDataMapper) SetTableName(table string) {
	t.table = table
}

func (t *TableDataMapper) TableName() string {
	return t.table
}

func (t *TableDataMapper) DBTableName() string {
	return t.db.TableName(t.table)
}
