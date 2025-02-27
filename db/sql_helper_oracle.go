package db

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-base/utils"
)

type OracleHelper struct{}

func init() {
	g_DBHelpers["oracle"] = &OracleHelper{}
}

func (helper *OracleHelper) InsertSqlBuilder(table string, params map[string]any, rawParams []string, options map[string]any) SqlBuildResult {
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
	paramCount := 1

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
			values += fmt.Sprintf(":p%d", paramCount)
			paralist = append(paralist, value)
			paramCount++
		} else {
			values += fmt.Sprintf("%v", value)
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, cols, values)
	return SqlBuildResult{Sql: sql, Params: paralist}
}

func (helper *OracleHelper) UpdateSQLBuilder(table string, params, filters map[string]any, wheres string) SqlBuildResult {
	var values string
	var filterstr string
	var wherestr string
	var paralist []any
	paramCount := 1

	params["DATA_UPD_DATE"] = utils.GetCurDate8()
	params["DATA_UPD_TIME"] = utils.GetCurDateTime14()

	for key, value := range params {
		if values != "" {
			values += ","
		}
		values += fmt.Sprintf("%s=:p%d", key, paramCount)
		paralist = append(paralist, value)
		paramCount++
	}

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
				filterstr += fmt.Sprintf(":p%d", paramCount)
				paralist = append(paralist, elem)
				paramCount++
			}
			filterstr += ")"
		default:
			filterstr += fmt.Sprintf("%s=:p%d", key, paramCount)
			paralist = append(paralist, value)
			paramCount++
		}
	}

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

func (helper *OracleHelper) QueryBaseSqlBuilder(sql string, params map[string]any) error {
	return nil
}

func (helper *OracleHelper) QueryNamedParamsBuilder(sqlOld string, params map[string]any) SqlBuildResult {
	var sqlParams []any
	sql := sqlOld
	paramCount := 1

	re := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`)
	paramIds := re.FindAllString(sql, -1)

	for _, p := range paramIds {
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
					varsBuilder.WriteString(fmt.Sprintf(":p%d", paramCount))
					sqlParams = append(sqlParams, paramValue.Index(i).Interface())
					paramCount++
				}
				sql = strings.Replace(sql, p, varsBuilder.String(), -1)
			default:
				sqlParams = append(sqlParams, param)
				sql = strings.Replace(sql, p, fmt.Sprintf(":p%d", paramCount), -1)
				paramCount++
			}
		} else {
			sql = strings.Replace(sql, p, "NULL", -1)
		}
	}

	return SqlBuildResult{Sql: sql, Params: sqlParams}
}

func (helper *OracleHelper) DeleteSQLBuilder(table string, filters map[string]any, wheres string) SqlBuildResult {
	var filterArr []string
	var paraList []any
	paramCount := 1

	for key, val := range filters {
		switch v := val.(type) {
		case []any:
			if len(v) > 0 {
				placeholders := make([]string, len(v))
				for i, item := range v {
					paraList = append(paraList, item)
					placeholders[i] = fmt.Sprintf(":p%d", paramCount)
					paramCount++
				}
				filterArr = append(filterArr, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ",")))
			}
		default:
			filterArr = append(filterArr, fmt.Sprintf("%s=:p%d", key, paramCount))
			paraList = append(paraList, val)
			paramCount++
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
