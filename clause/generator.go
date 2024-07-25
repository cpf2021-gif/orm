package clause

import (
	"fmt"
	"strings"
)

type generator func(vals ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = insertClause
	generators[VALUES] = valuesClause
	generators[SELECT] = selectClause
	generators[LIMIT] = limitClause
	generators[WHERE] = whereClause
	generators[ORDERBY] = orderbyClause
	generators[UPDATE] = updateClause
	generators[DELETE] = deleteClause
	generators[COUNT] = countClause
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func insertClause(vals ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := vals[0]
	fileds := strings.Join(vals[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %v (%s)", tableName, fileds), []interface{}{}
}

func valuesClause(vals ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bingStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range vals {
		v := value.([]interface{})
		if bingStr == "" {
			bingStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%s)", bingStr))
		if i != len(vals)-1 {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func selectClause(vals ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := vals[0]
	filelds := strings.Join(vals[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %v", filelds, tableName), []interface{}{}
}

func limitClause(vals ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", vals
}

func whereClause(vals ...interface{}) (string, []interface{}) {
	// WHERE $exp
	desc, vars := vals[0], vals[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func orderbyClause(vals ...interface{}) (string, []interface{}) {
	// ORDER BY
	return fmt.Sprintf("ORDER BY %s", vals[0]), []interface{}{}
}

func updateClause(vals ...interface{}) (string, []interface{}) {
	tableName := vals[0]
	m := vals[1].(map[string]interface{})

	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %v SET %s", tableName, strings.Join(keys, ", ")), vars
}

func deleteClause(vals ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %v", vals[0]), []interface{}{}
}

func countClause(vals ...interface{}) (string, []interface{}) {
	return selectClause(vals[0], []string{"count(*)"})
}
