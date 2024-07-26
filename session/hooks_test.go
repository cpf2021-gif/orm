package session

import (
	"testing"

	"orm/log"
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", a)
	return nil
}

func (a *Account) AfterInsert(s *Session) error {
	log.Info("after insert")
	return nil
}

func (a *Account) BeforeQuery(s *Session) error {
	log.Info("before query")
	return nil
}

func (a *Account) AfterQuery(s *Session) error {
	log.Info("after query", a)
	a.Password = "******"
	return nil
}

func TestSessionCallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.Password != "******" || u.ID != 1 {
		t.Fatal("failed to call hooks after query, got", u.Password)
	}
}
