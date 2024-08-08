package db

import (
	"fmt"
	"testing"
)

func TestInsertSqlBuilder(t *testing.T) {
	table := "your_table"
	params := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	rawParams := []string{}
	options := map[string]interface{}{
		//"defaultDate": true,
	}

	result := InsertSqlBuilder(table, params, rawParams, options)
	fmt.Println("InsertSqlBuilder result:")
	fmt.Println(result.Sql)
	fmt.Println(result.Params)
}

func TestUpdateSqlBuilder(t *testing.T) {
	table := "my_table"
	params := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	filters := map[string]interface{}{
		"status": []interface{}{"active", "pending"},
	}
	wheres := "id=1"

	result := UpdateSQLBuilder(table, params, filters, wheres)
	fmt.Println("UpdateSQLBuilder result:")
	fmt.Println(result.Sql)
	fmt.Println(result.Params)
}

func TestDeleteSQLBuilder(t *testing.T) {
	table := "users"
	filters := map[string]interface{}{
		"id":   []interface{}{1, 2, 3},
		"name": "John Doe",
	}
	//wheres := "age > 30"
	wheres := ""

	result := DeleteSQLBuilder(table, filters, wheres)
	fmt.Println("DeleteSQLBuilder result:")
	fmt.Println(result.Sql)
	fmt.Println(result.Params)
}

func TestQueryBaseSqlBuilder(t *testing.T) {

}

func TestQueryNamedParamsBuilder(t *testing.T) {
	sqlOld := "SELECT * FROM users WHERE id = {user_id} AND status IN ({status_list})"
	params := map[string]interface{}{
		"USER_ID":     123,
		"STATUS_LIST": []interface{}{"active", "pending"},
	}

	result := QueryNamedParamsBuilder(sqlOld, params)
	fmt.Println("TestQueryNamedParamsBuilder result:")
	fmt.Println(result.Sql)
	fmt.Println(result.Params)
}

// func TestMain(m *testing.M) {

// }
