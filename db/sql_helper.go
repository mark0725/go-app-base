package db

import (
	"fmt"
)

type ISqlHelper interface {
	InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult
	UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult
	QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult
	DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult
	QueryBaseSqlBuilder(sql string, params map[string]any) error
}

var g_DBHelpers = map[string]ISqlHelper{}

func InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult {
	conn := GetDBConn(DB_CONN_NAME_DEFAULT)
	if conn == nil {
		return SqlBuildResult{}
	}

	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return SqlBuildResult{}
	}

	return helper.InsertSqlBuilder(table, params, rawParams, options)
}

func UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult {
	conn := GetDBConn(DB_CONN_NAME_DEFAULT)
	if conn == nil {
		return SqlBuildResult{}
	}
	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return SqlBuildResult{}
	}
	return helper.UpdateSQLBuilder(table, params, filters, wheres)
}

func QueryBaseSqlBuilder(sql string, params map[string]any) error {
	conn := GetDBConn(DB_CONN_NAME_DEFAULT)
	if conn == nil {
		return fmt.Errorf("not found db: %s", conn.Name)
	}
	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return fmt.Errorf("not found db: %s", conn.Name)
	}
	return helper.QueryBaseSqlBuilder(sql, params)
}

// QueryNamedParamsBuilder function
func QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult {
	conn := GetDBConn(DB_CONN_NAME_DEFAULT)
	if conn == nil {
		return SqlBuildResult{}
	}
	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return SqlBuildResult{}
	}

	return helper.QueryNamedParamsBuilder(sqlOld, params)
}

func DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult {
	conn := GetDBConn(DB_CONN_NAME_DEFAULT)
	if conn == nil {
		return SqlBuildResult{}
	}
	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return SqlBuildResult{}
	}
	return helper.DeleteSQLBuilder(table, filters, wheres)
}

func GetSqlHelper(db string) ISqlHelper {
	conn := GetDBConn(db)
	if conn == nil {
		return nil
	}

	if conn.Type == "" {
		return nil
	}

	if _, exists := g_DBHelpers[conn.Type]; !exists {
		return nil
	}

	return g_DBHelpers[conn.Type]
}

func RegisterSqlHelper(dbType string, helper ISqlHelper) {
	g_DBHelpers[dbType] = helper
}
