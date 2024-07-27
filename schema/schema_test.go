package schema

import (
	"testing"

	"orm/dialect"
)

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func (u *User) TableName() string {
	return "users"
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "users" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to pase primary key")
	}
}
