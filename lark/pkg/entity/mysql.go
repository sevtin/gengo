package entity

import (
	"fmt"
	"reflect"
	"time"
)

type GormCreatedTs struct {
	CreatedTs int64 `gorm:"column:created_ts;autoCreateTime:milli;NOT NULL" json:"created_ts"`
}

type GormUpdatedTs struct {
	UpdatedTs int64 `gorm:"column:updated_ts;autoUpdateTime:milli;NOT NULL" json:"updated_ts"`
}

type GormDeletedTs struct {
	DeletedTs int64 `gorm:"column:deleted_ts;default:0;NOT NULL" json:"deleted_ts"`
}

type GormEntityTs struct {
	CreatedTs int64 `gorm:"column:created_ts;autoCreateTime:milli;NOT NULL" json:"created_ts"`
	UpdatedTs int64 `gorm:"column:updated_ts;autoUpdateTime:milli;NOT NULL" json:"updated_ts"`
	DeletedTs int64 `gorm:"column:deleted_ts;default:0;NOT NULL" json:"deleted_ts"`
}

type GormTs struct {
	CreatedTs int64 `gorm:"column:created_ts;autoCreateTime:milli;NOT NULL" json:"created_ts"`
	UpdatedTs int64 `gorm:"column:updated_ts;autoUpdateTime:milli;NOT NULL" json:"updated_ts"`
}

func Deleted() (column string, value interface{}) {
	return "deleted_ts", time.Now().UnixMilli()
}

type Mysql struct {
	Query string
	Args  []interface{}
}

func (m *Mysql) andCondition(condition string) {
	if m.Query == "" {
		m.Query = condition
	} else {
		m.Query += " AND " + condition
	}
}

func (m *Mysql) orCondition(condition string) {
	if m.Query == "" {
		m.Query = condition
	} else {
		m.Query = "(" + m.Query + ") OR (" + condition + ")"
	}
}

func (m *Mysql) NotDeleted(alias string) {
	clause := alias + ".deleted_ts=0"
	m.andCondition(clause)
}

func (m *Mysql) IsNull(field string) {
	clause := field + " IS NULL"
	m.andCondition(clause)
}

func (m *Mysql) IsNotNull(field string) {
	clause := field + " IS NOT NULL"
	m.andCondition(clause)
}

func (m *Mysql) Where(query string, value ...interface{}) {
	m.andCondition(query)
	m.Args = append(m.Args, value...)
}

func (m *Mysql) OrWhere(query string, value ...interface{}) {
	m.orCondition(query)
	m.Args = append(m.Args, value...)
}

func (m *Mysql) Between(field string, begin interface{}, end interface{}) {
	clause := field + " BETWEEN ? AND ?"
	m.andCondition(clause)
	m.Args = append(m.Args, begin)
	m.Args = append(m.Args, end)
}

func (m *Mysql) OpenParen() {
	m.Query += "(1=1"
}

func (m *Mysql) CloseParen() {
	m.Query += ")"
}

func (m *Mysql) AndQuery(query string) {
	m.andCondition(query)
}

func (m *Mysql) OrQuery(query string) {
	m.orCondition(query)
}

func (m *Mysql) SetLink(link string) {
	m.Query += fmt.Sprintf(" %s ", link)
}

func (m *Mysql) AppendArg(value interface{}) {
	m.Args = append(m.Args, value)
}

func (m *Mysql) SetConditions(obj any) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		tag := field.Tag.Get("where")
		if tag != "" {
			m.Where(tag+"=?", value)
		}
	}
}
