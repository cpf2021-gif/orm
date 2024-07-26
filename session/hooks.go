package session

import "reflect"

type MethodType int

const (
	BeforeQuery MethodType = iota
	AfterQuery

	BeforeUpdate
	AfterUpdate

	BeforeDelete
	AfterDelete

	BeforeInsert
	AfterInsert
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}
type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}
type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}

type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}
type IAfterDelete interface {
	AfterDelete(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}
type IAfterInsert interface {
	AfterInsert(s *Session) error
}

func (s *Session) CallMethod(method MethodType, value interface{}) {
	dest := reflect.ValueOf(value)
	// 获取方法
	switch method {
	case BeforeQuery:
		if v, ok := dest.Interface().(IBeforeQuery); ok {
			v.BeforeQuery(s)
		}
	case AfterQuery:
		if v, ok := dest.Interface().(IAfterQuery); ok {
			v.AfterQuery(s)
		}
	case BeforeUpdate:
		if v, ok := dest.Interface().(IBeforeUpdate); ok {
			v.BeforeUpdate(s)
		}
	case AfterUpdate:
		if v, ok := dest.Interface().(IAfterUpdate); ok {
			v.AfterUpdate(s)
		}
	case BeforeDelete:
		if v, ok := dest.Interface().(IBeforeDelete); ok {
			v.BeforeDelete(s)
		}
	case AfterDelete:
		if v, ok := dest.Interface().(IAfterDelete); ok {
			v.AfterDelete(s)
		}
	case BeforeInsert:
		if v, ok := dest.Interface().(IBeforeInsert); ok {
			v.BeforeInsert(s)
		}
	case AfterInsert:
		if v, ok := dest.Interface().(IAfterInsert); ok {
			v.AfterInsert(s)
		}
	}
}
