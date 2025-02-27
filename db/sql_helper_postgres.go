package db

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-base/utils"
)

// PostgreSQL helper struct
type PostgresHelper struct {
}

func init() {
	g_DBHelpers["postgres"] = &PostgresHelper{}
}

// InsertSqlBuilder for PostgreSQL
func (helper *PostgresHelper) InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult {
	defaultDate := true
	if options != nil {
		if val, ok := options["defaultDate"].(bool); ok {
			defaultDate = val
		}
	}

	cols := ""
	values := ""
	var paralist []any

	paras := make(map[string]any)
	for k, v := range params {
		paras[k] = v
	}
	if defaultDate {
		paras["data_crt_date"] = utils.GetCurDate8()
		paras["data_crt_time"] = utils.GetCurDateTime14()
		paras["data_upd_date"] = utils.GetCurDate8()
		paras["data_upd_time"] = utils.GetCurDateTime14()
	}

	for key, value := range paras {
		if cols != "" {
			cols += ","
			values += ","
		}
		cols += key
		isRaw := false
		for _, rawParam := range rawParams {
			if rawParam == key {
				isRaw = true
				break
			}
		}
		if !isRaw {
			values += fmt.Sprintf("$%d", len(paralist)+1) // PostgreSQL uses $n placeholders
			paralist = append(paralist, value)
		} else {
			values += fmt.Sprintf("%v", value)
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, cols, values)
	return SqlBuildResult{Sql: sql, Params: paralist}
}

// UpdateSQLBuilder for PostgreSQL
func (helper *PostgresHelper) UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult {
	var values string
	var filterstr string
	var wherestr string
	var paralist []any

	// Merge params with update timestamps
	params["data_upd_date"] = utils.GetCurDate8()
	params["data_upd_time"] = utils.GetCurDateTime14()

	// Build the SET clause
	for key, value := range params {
		if values != "" {
			values += ","
		}
		values += fmt.Sprintf("%s=$%d", key, len(paralist)+1)
		paralist = append(paralist, value)
	}

	// Build the WHERE clause with filters
	for key, value := range filters {
		if filterstr != "" {
			filterstr += " AND "
		}
		switch v := value.(type) {
		case []any:
			filterstr += fmt.Sprintf("%s = ANY ($%d)", key, len(paralist)+1) // Support for array/vector fields
			paralist = append(paralist, v)
		default:
			filterstr += fmt.Sprintf("%s=$%d", key, len(paralist)+1)
			paralist = append(paralist, value)
		}
	}

	// Combine custom WHERE clause and filters
	wherestr = filterstr
	if wheres != "" {
		if filterstr == "" {
			wherestr = wheres
		} else {
			wherestr = fmt.Sprintf("%s AND %s", wheres, filterstr)
		}
	}
	if wherestr != "" {
		wherestr = " WHERE " + wherestr
	}

	sql := fmt.Sprintf("UPDATE %s SET %s%s", table, values, wherestr)
	return SqlBuildResult{Sql: sql, Params: paralist}
}

// QueryBaseSqlBuilder (placeholder, as in the original)
func (helper *PostgresHelper) QueryBaseSqlBuilder(sql string, params map[string]any) error {
	return nil
}

// QueryNamedParamsBuilder for PostgreSQL
func (helper *PostgresHelper) QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult {
	var sqlParams []any
	sql := sqlOld

	re := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`)
	paramIds := re.FindAllString(sql, -1)

	for _, p := range paramIds {
		varName := strings.ToUpper(p[1 : len(p)-1])

		if param, ok := params[varName]; ok {
			paramValue := reflect.ValueOf(param)
			switch paramValue.Kind() {
			case reflect.Slice, reflect.Array:
				sql = strings.Replace(sql, p, fmt.Sprintf("$%d", len(sqlParams)+1), -1)
				sqlParams = append(sqlParams, param)
			default:
				sql = strings.Replace(sql, p, fmt.Sprintf("$%d", len(sqlParams)+1), -1)
				sqlParams = append(sqlParams, param)
			}
		} else {
			sql = strings.Replace(sql, p, "NULL", -1)
		}
	}

	return SqlBuildResult{Sql: sql, Params: sqlParams}
}

// DeleteSQLBuilder for PostgreSQL
func (helper *PostgresHelper) DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult {
	var filterArr []string
	var paraList []any

	for key, val := range filters {
		switch v := val.(type) {
		case []any:
			if len(v) > 0 {
				filterArr = append(filterArr, fmt.Sprintf("%s = ANY ($%d)", key, len(paraList)+1))
				paraList = append(paraList, v) // Support for vector/array fields
			}
		default:
			filterArr = append(filterArr, fmt.Sprintf("%s=$%d", key, len(paraList)+1))
			paraList = append(paraList, val)
		}
	}

	if wheres != "" {
		filterArr = append([]string{wheres}, filterArr...)
	}

	whereStr := ""
	if len(filterArr) > 0 {
		whereStr = " WHERE " + strings.Join(filterArr, " AND ")
	}

	sql := fmt.Sprintf("DELETE FROM %s%s", table, whereStr)
	return SqlBuildResult{Sql: sql, Params: paraList}
}
