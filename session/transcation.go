package session

import "orm/log"

func (s *Session) Begin() (err error) {
	log.Info("begin transaction")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("commit transaction")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("rollback transaction")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
