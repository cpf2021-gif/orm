package orm

import (
	"errors"
	"testing"

	"orm/session"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "orm.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func (u *User) TableName() string {
	return "users"
}

func TestEngineTranscation(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})

	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, _ = s.Insert(&User{Name: "Tom", Age: 18})
		return nil, errors.New("Error")
	})

	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, _ = s.Insert(&User{Name: "Tom", Age: 18})
		return nil, nil
	})

	count, _ := s.Model(&User{}).Count()

	if err != nil || !s.HasTable() || count != 1 {
		t.Fatal("failed to commit")
	}
}

func TestEngine(t *testing.T) {
	t.Run("alias table name", func(t *testing.T) {
		aliasTableNames(t)
	})
}

func aliasTableNames(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{})
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()

	if s.RefTable().Name != "users" || !s.HasTable() {
		t.Fatal("failed to alias table name")
	}
}
