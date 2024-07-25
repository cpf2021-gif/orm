package session

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"orm/log"
	"orm/schema"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	if table == nil {
		return errors.New("model is not set")
	}
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	table := s.RefTable()
	if table == nil {
		return errors.New("model is not set")
	}

	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	table := s.RefTable()
	if table == nil {
		log.Error("model is not set")
		return false
	}

	sql, values := s.dialect.TableExistSQL(table.Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == table.Name
}
