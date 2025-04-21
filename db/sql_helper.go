package db

import (
	"fmt"
	"strings"
)

type ISqlHelper interface {
	InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult
	UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult
	QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult
	DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult
	QueryBaseSqlBuilder(sql string, params map[string]any) error
	PageQueryBuilder(sql string, limit int, offset int) (string, error)
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

func PageQueryBuilder(db string, sql string, limit int, offset int) (string, error) {
	conn := GetDBConn(db)
	if conn == nil {
		return "", fmt.Errorf("not found db: %s", db)
	}

	helper := g_DBHelpers[conn.Type]
	if helper == nil {
		return "", fmt.Errorf("not found db helper: %s", conn.Name)
	}

	return helper.PageQueryBuilder(sql, limit, offset)
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

func QueryParamsBuilder(params map[string]any, defaultPrefix string, otherPrefixes map[string][]string, spOps map[string]string) string {
	sqlParts := []string{}
	paramMap := make(map[string]any)
	//paramCount := make(map[string]int)

	getPlaceholder := func(key string) string {
		//paramCount[key]++
		// return fmt.Sprintf("{%s_%d}", key, paramCount[key])
		return fmt.Sprintf("{%s}", key)
	}

	for k, val := range params {
		paraName := k
		if defaultPrefix != "" {
			paraName = defaultPrefix + "." + k
		}

		for p, fields := range otherPrefixes {
			for _, f := range fields {
				if f == k {
					paraName = p + "." + k
					break
				}
			}
		}

		spOp := ""
		if spOps != nil {
			spOp = spOps[k]
		}

		buildLike := func(arr []any, single bool) string {
			if single {
				strVal, _ := arr[0].(string)
				strVal = strings.ReplaceAll(strVal, "'", "")
				return fmt.Sprintf("%s like '%%%s%%'", paraName, strVal)
			}
			partArr := []string{}
			for i, v := range arr {
				strVal, _ := v.(string)
				strVal = strings.ReplaceAll(strVal, "'", "")
				partArr = append(partArr, fmt.Sprintf("%s like '%%%s%%'", paraName, strVal))
				arr[i] = strVal
			}
			return fmt.Sprintf("(%s)", strings.Join(partArr, " or "))
		}

		switch spOp {
		case "like":
			if arrVal, ok := val.([]any); ok && len(arrVal) > 0 {
				sqlParts = append(sqlParts, buildLike(arrVal, false))
			} else {
				strVal, _ := val.(string)
				valArr := []any{strVal}
				sqlParts = append(sqlParts, buildLike(valArr, true))
			}
		case "-":
			// 区间
			if arrVal, ok := val.([]any); ok {
				if len(arrVal) < 2 {
					// paraName <= ?
					strVal, _ := arrVal[0].(string)
					strVal = strings.ReplaceAll(strVal, "-", "")
					ph := getPlaceholder(k)
					sqlParts = append(sqlParts, fmt.Sprintf("%s <= %s", paraName, ph))
					paramMap[ph] = strVal
				} else {
					// (paraName >= ? and paraName <= ?)
					strVal0, _ := arrVal[0].(string)
					strVal0 = strings.ReplaceAll(strVal0, "-", "")
					strVal1, _ := arrVal[1].(string)
					strVal1 = strings.ReplaceAll(strVal1, "-", "")

					ph0 := getPlaceholder(k)
					ph1 := getPlaceholder(k)
					sqlParts = append(sqlParts, fmt.Sprintf("(%s >= %s and %s <= %s)", paraName, ph0, paraName, ph1))
					paramMap[ph0] = strVal0
					paramMap[ph1] = strVal1
				}
			} else {
				strVal, _ := val.(string)
				if !strings.Contains(strVal, "-") {
					ph := getPlaceholder(k)
					sqlParts = append(sqlParts, fmt.Sprintf("%s <= %s", paraName, ph))
					paramMap[ph] = strings.ReplaceAll(strVal, "-", "")
				}
			}
		case "not in":
			// not in
			if arrVal, ok := val.([]any); ok {
				placeholders := []string{}
				for _, one := range arrVal {
					ph := getPlaceholder(k)
					placeholders = append(placeholders, ph)
					paramMap[ph] = one
				}
				sqlParts = append(sqlParts, fmt.Sprintf("%s not in(%s)", paraName, strings.Join(placeholders, ",")))
			}
		default:
			// no special op or other operators
			if spOp != "" && spOp != "-" && spOp != "not in" && spOp != "like" {
				// paraName + spOps[k] + ?
				ph := getPlaceholder(k)
				sqlParts = append(sqlParts, fmt.Sprintf("%s %s %s", paraName, spOp, ph))
				paramMap[ph] = val
			} else {
				// in or default =
				if arrVal, ok := val.([]any); ok {
					placeholders := []string{}
					for _, one := range arrVal {
						ph := getPlaceholder(k)
						placeholders = append(placeholders, ph)
						paramMap[ph] = one
					}
					sqlParts = append(sqlParts, fmt.Sprintf("%s in(%s)", paraName, strings.Join(placeholders, ",")))
				} else {
					ph := getPlaceholder(k)
					sqlParts = append(sqlParts, fmt.Sprintf("%s = %s", paraName, ph))
					paramMap[ph] = val
				}
			}
		}
	}

	return strings.Join(sqlParts, " and ")
}

func RegisterSqlHelper(dbType string, helper ISqlHelper) {
	g_DBHelpers[dbType] = helper
}
