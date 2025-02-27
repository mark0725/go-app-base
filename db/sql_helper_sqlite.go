package db

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-base/utils"
)

type SqliteHelper struct{}

func init() {
	g_DBHelpers["sqlite"] = &SqliteHelper{}
}

func (helper *SqliteHelper) InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult {
	defaultDate := true
	if options != nil {
		defaultDate, ok := options["defaultDate"].(bool)
		if !ok || defaultDate {
			defaultDate = true
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
		paras["DATA_CRT_DATE"] = utils.GetCurDate8()
		paras["DATA_CRT_TIME"] = utils.GetCurDateTime14()
		paras["DATA_UPD_DATE"] = utils.GetCurDate8()
		paras["DATA_UPD_TIME"] = utils.GetCurDateTime14()
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
			values += "?"
			paralist = append(paralist, value)
		} else {
			values += fmt.Sprintf("%v", value)
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, cols, values)

	return SqlBuildResult{Sql: sql, Params: paralist}
}

func (helper *SqliteHelper) UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult {
	var values string
	var filterstr string
	var wherestr string
	var paralist []any

	// Merge params with DATA_UPD_DATE and DATA_UPD_TIME
	params["DATA_UPD_DATE"] = utils.GetCurDate8()
	params["DATA_UPD_TIME"] = utils.GetCurDateTime14()

	// Build the values part
	for key, value := range params {
		if values != "" {
			values += ","
		}
		values += fmt.Sprintf("%s=?", key)
		paralist = append(paralist, value)
	}

	// Build the filters part
	for key, value := range filters {
		if filterstr != "" {
			filterstr += " AND "
		}
		switch v := value.(type) {
		case []any:
			filterstr += fmt.Sprintf("%s IN (", key)
			for i, elem := range v {
				if i > 0 {
					filterstr += ","
				}
				filterstr += "?"
				paralist = append(paralist, elem)
			}
			filterstr += ")"
		default:
			filterstr += fmt.Sprintf("%s=?", key)
			paralist = append(paralist, value)
		}
	}

	// Combine wheres and filterstr
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

func (helper *SqliteHelper) QueryBaseSqlBuilder(sql string, params map[string]any) error {
	return nil
}

func (helper *SqliteHelper) QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult {
	var sqlParams []any
	sql := sqlOld

	re := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`)
	paramIds := re.FindAllString(sql, -1)

	for _, p := range paramIds {
		vars := "?"
		varName := strings.ToUpper(p[1 : len(p)-1])

		if param, ok := params[varName]; ok {
			paramValue := reflect.ValueOf(param)
			switch paramValue.Kind() {
			case reflect.Slice, reflect.Array:
				var varsBuilder strings.Builder
				for i := 0; i < paramValue.Len(); i++ {
					if i > 0 {
						varsBuilder.WriteString(",")
					}
					sqlParams = append(sqlParams, paramValue.Index(i).Interface())
					varsBuilder.WriteString("?")
				}
				sql = strings.Replace(sql, p, varsBuilder.String(), -1)
			default:
				sqlParams = append(sqlParams, param)
				sql = strings.Replace(sql, p, vars, -1)
			}
		} else {
			sql = strings.Replace(sql, p, "NULL", -1)
		}
	}

	return SqlBuildResult{Sql: sql, Params: sqlParams}
}

func (helper *SqliteHelper) DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult {
	var filterArr []string
	var paraList []any

	for key, val := range filters {
		switch v := val.(type) {
		case []any:
			if len(v) > 0 {
				placeholders := make([]string, len(v))
				for i, item := range v {
					paraList = append(paraList, item)
					placeholders[i] = "?"
				}
				filterArr = append(filterArr, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ",")))
			}
		default:
			filterArr = append(filterArr, fmt.Sprintf("%s=?", key))
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
