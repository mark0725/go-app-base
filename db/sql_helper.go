package db

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-base/utils"
)

type SqlBuildResult struct {
	Sql    string
	Params []interface{}
}

func InsertSqlBuilder(table string, params map[string]interface{}, rawParams []string, options map[string]interface{}) SqlBuildResult {
	defaultDate := true
	if options != nil {
		defaultDate, ok := options["defaultDate"].(bool)
		if !ok || defaultDate {
			defaultDate = true
		}
	}

	cols := ""
	values := ""
	var paralist []interface{}

	paras := make(map[string]interface{})

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

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", table, cols, values)

	return SqlBuildResult{Sql: sql, Params: paralist}

}

func UpdateSQLBuilder(table string, params, filters map[string]interface{}, wheres string) SqlBuildResult {
	var values string
	var filterstr string
	var wherestr string
	var paralist []interface{}

	// Merge params with DATA_UPD_DATE and DATA_UPD_TIME.
	params["DATA_UPD_DATE"] = utils.GetCurDate8()
	params["DATA_UPD_TIME"] = utils.GetCurDateTime14()

	// Build the values part of the SQL statement.
	for key, value := range params {
		if values != "" {
			values += ","
		}
		values += fmt.Sprintf("%s=?", key)
		paralist = append(paralist, value)
	}

	// Build the filters part of the SQL statement.
	for key, value := range filters {
		if filterstr != "" {
			filterstr += " and "
		}
		switch v := value.(type) {
		case []interface{}:
			filterstr += fmt.Sprintf("%s in (", key)
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

	// Combine wheres and filterstr to form the full where clause.
	wherestr = filterstr
	if wheres != "" {
		if filterstr == "" {
			wherestr = wheres
		} else {
			wherestr = fmt.Sprintf("%s and %s", wheres, filterstr)
		}
	}
	if wherestr != "" {
		wherestr = " where " + wherestr
	}

	// Build the final SQL statement.
	sql := fmt.Sprintf("update %s set %s%s", table, values, wherestr)
	return SqlBuildResult{Sql: sql, Params: paralist}
}

func QueryBaseSqlBuilder(sql string, params map[string]interface{}) error {
	return nil

}

// QueryNamedParamsBuilder function
func QueryNamedParamsBuilder(sqlOld string, params map[string]interface{}) SqlBuildResult {
	var sqlParams []interface{}
	sql := sqlOld

	re := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`)
	paramIds := re.FindAllString(sql, -1)

	for _, p := range paramIds {
		vars := "?"
		varName := strings.ToUpper(p[1 : len(p)-1])

		if param, ok := params[varName]; ok {
			paramValue := reflect.ValueOf(param)
			//fmt.Printf("param %s type: %v\n", varName, paramValue.Kind())
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
			sql = strings.Replace(sql, p, "null", -1)
		}
	}

	return SqlBuildResult{Sql: sql, Params: sqlParams}
}

func DeleteSQLBuilder(table string, filters map[string]interface{}, wheres string) SqlBuildResult {
	var filterArr []string
	var paraList []interface{}

	for key, val := range filters {
		switch v := val.(type) {
		case []interface{}:
			if len(v) > 0 {
				placeholders := make([]string, len(v))
				for i, item := range v {
					paraList = append(paraList, item)
					placeholders[i] = "?"
				}
				filterArr = append(filterArr, fmt.Sprintf("%s in (%s)", key, strings.Join(placeholders, ",")))
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
		whereStr = " where " + strings.Join(filterArr, " and ")
	}

	sql := fmt.Sprintf("delete from %s%s", table, whereStr)

	return SqlBuildResult{Sql: sql, Params: paraList}
}
