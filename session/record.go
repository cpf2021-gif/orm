package session

import (
	"errors"
	"reflect"

	"orm/clause"
	"orm/schema"
)

// Insert(&User{}) or Insert([]&User{})
func (s *Session) Insert(val interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	reflectValue := reflect.ValueOf(val)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	var table *schema.Schema

	if reflectValue.Kind() == reflect.Slice {
		for i := 0; i < reflectValue.Len(); i++ {
			if i == 0 {
				table = s.Model(reflectValue.Index(i).Interface()).RefTable()
			}
			s.CallMethod(BeforeInsert, reflectValue.Index(i).Interface())
			recordValues = append(recordValues, table.RecordValues(reflectValue.Index(i).Interface()))
		}
	} else if reflectValue.Kind() == reflect.Struct {
		table = s.Model(val).RefTable()
		s.CallMethod(BeforeInsert, val)
		recordValues = append(recordValues, table.RecordValues(val))
	} else {
		return 0, errors.New("unsupported type")
	}

	s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	if reflectValue.Kind() == reflect.Slice {
		for i := 0; i < reflectValue.Len(); i++ {
			s.CallMethod(AfterInsert, reflectValue.Index(i).Interface())
		}
	} else {
		s.CallMethod(AfterInsert, val)
	}
	return result.RowsAffected()
}

func (s *Session) Find(vals interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(vals))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.CallMethod(BeforeQuery, reflect.New(destType).Elem().Addr().Interface())

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		if len(kv)%2 != 0 {
			return 0, errors.New("kv invaild")
		}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	table := s.RefTable()
	if table == nil || table.Name == "" {
		return 0, errors.New("no set model")
	}

	s.clause.Set(clause.UPDATE, table.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	table := s.RefTable()
	if table == nil || table.Name == "" {
		return 0, errors.New("no set model")
	}

	s.clause.Set(clause.DELETE, table.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	table := s.RefTable()
	if table == nil || table.Name == "" {
		return 0, errors.New("no set model")
	}

	s.clause.Set(clause.COUNT, table.Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
