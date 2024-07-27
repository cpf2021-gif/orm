package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"orm/dialect"
	"orm/log"
	"orm/session"
)

type Engine struct {
	db      *sql.DB
	dislect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dislect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close detabase")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dislect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.Rollback()
		} else {
			err = s.Commit()
			if err != nil {
				s.Rollback()
			}
		}
	}()

	return f(s)
}

func difference(a, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}

	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Info("Table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		if table == nil {
			return nil, errors.New("model is not a struct")
		}
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		err = rows.Close()
		if err != nil {
			return
		}

		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("Added cols: %v, Deleted cols: %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			_, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)).Exec()
			if err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
