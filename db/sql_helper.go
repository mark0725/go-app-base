package db

import (
	"fmt"
	"strings"

	"github.com/mark0725/go-app-base/utils"
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

func QueryNamedParamsBuilder2(db string, sqlOld string, params map[string]any) SqlBuildResult {
	conn := GetDBConn(db)
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

type TYPE_WHERECLAUSE string

const (
	TYPE_WHERECLAUSE_STRING TYPE_WHERECLAUSE = "string"
)

type whereClause struct {
	whereType TYPE_WHERECLAUSE
	query     string
	params    map[string]any
	queryMap  map[string]any
	fields    []string
}

type QueryParamsOptions struct {
	sql           string
	cols          []string
	from          string
	where         *whereClause
	params        map[string]any
	defaultPrefix string
	otherPrefixes map[string][]string
	spOps         map[string]string
	groupBy       string
	orderBy       string
	limit         int
	offset        int
}

func NewQueryParamsBuilder() *QueryParamsOptions {
	return &QueryParamsOptions{
		sql:           "",
		cols:          []string{},
		from:          "",
		params:        make(map[string]any),
		defaultPrefix: "",
		otherPrefixes: make(map[string][]string),
		spOps:         make(map[string]string),
		groupBy:       "",
		orderBy:       "",
	}
}
func (q *QueryParamsOptions) DefaultPrefix(prefix string) *QueryParamsOptions {
	q.defaultPrefix = prefix
	return q
}
func (q *QueryParamsOptions) OtherPrefixes(otherPrefixes map[string][]string) *QueryParamsOptions {
	q.otherPrefixes = otherPrefixes
	return q
}

func (q *QueryParamsOptions) SpOps(spOps map[string]string) *QueryParamsOptions {
	q.spOps = spOps
	return q
}

func (q *QueryParamsOptions) Sql(sql string) *QueryParamsOptions {
	q.sql = sql
	return q
}

func (q *QueryParamsOptions) Columns(cols ...string) *QueryParamsOptions {
	q.cols = cols
	return q
}

func (q *QueryParamsOptions) From(from string) *QueryParamsOptions {
	q.from = from
	return q
}

func (q *QueryParamsOptions) WhereString(query string, params map[string]any) *QueryParamsOptions {

	q.where = &whereClause{
		whereType: TYPE_WHERECLAUSE_STRING,
		query:     query,
		params:    params,
	}
	return q
}

func (q *QueryParamsOptions) WhereMap(params map[string]any, field ...string) *QueryParamsOptions {

	for k, v := range params {
		if len(field) > 0 && !utils.Contains(field, k) {
			continue
		}
		q.params[k] = v
	}

	return q
}

func (q *QueryParamsOptions) Group(groupBy string) *QueryParamsOptions {
	q.groupBy = groupBy
	return q
}

func (q *QueryParamsOptions) Order(orderBy string) *QueryParamsOptions {
	q.orderBy = orderBy
	return q
}

func (q *QueryParamsOptions) Limit(limit int) *QueryParamsOptions {
	q.limit = limit
	return q
}

func (q *QueryParamsOptions) Offset(offset int) *QueryParamsOptions {
	q.offset = offset
	return q
}

func (q *QueryParamsOptions) Build() (string, map[string]any) {
	cols := "*"
	if len(q.cols) > 0 {
		cols = strings.Join(q.cols, ",")
	}
	sql := fmt.Sprintf("select %s from %s", cols, q.from)
	sqlParams := map[string]any{}
	if q.where != nil {
		sql += " where " + q.where.query
		for k, v := range q.where.params {
			sqlParams[k] = v
		}
	} else {
		sql += " where 1=1"
	}

	if q.sql != "" {
		sql += q.sql
	}

	sqlParts, sqlPartParams := q.BuildSqlPart()
	if sqlParts != "" {
		sql += " and (" + sqlParts + ")"
	}
	if q.groupBy != "" {
		sql += " group by " + q.groupBy
	}
	if q.orderBy != "" {
		sql += " order by " + q.orderBy
	}
	for k, v := range sqlPartParams {
		sqlParams[k] = v
	}
	return sql, sqlParams
}

func (q *QueryParamsOptions) BuildSqlPart() (string, map[string]any) {
	sqlParts := []string{}
	paramMap := make(map[string]any)
	paramCount := make(map[string]int)

	getPlaceholder := func(key string, v any) string {
		paramCount[key]++
		paramKey := fmt.Sprintf("%s_%d", key, paramCount[key])
		paramMap[paramKey] = v
		return fmt.Sprintf("{%s}", paramKey)
	}

	for k, val := range q.params {
		paraName := k
		if q.defaultPrefix != "" {
			paraName = q.defaultPrefix + "." + k
		}

		for p, fields := range q.otherPrefixes {
			for _, f := range fields {
				if f == k {
					paraName = p + "." + k
					break
				}
			}
		}

		spOp := ""
		if q.spOps != nil && q.spOps[k] != "" {
			spOp = q.spOps[k]
		}

		switch spOp {
		case "like":
			if arrVal, ok := val.([]any); ok && len(arrVal) > 0 {
				partArr := []string{}
				for _, v := range arrVal {
					strVal, _ := v.(string)
					strVal = strings.ReplaceAll(strVal, "'", "")
					ph := getPlaceholder(k, fmt.Sprintf("%%%s%%", strVal))
					partArr = append(partArr, fmt.Sprintf("%s like %s", paraName, ph))
				}

				sqlParts = append(sqlParts, fmt.Sprintf("(%s)", strings.Join(partArr, " or ")))
			} else {
				strVal, _ := val.(string)
				strVal = strings.ReplaceAll(strVal, "'", "")
				ph := getPlaceholder(k, fmt.Sprintf("%%%s%%", strVal))
				sqlParts = append(sqlParts, fmt.Sprintf("%s like %s", paraName, ph))
			}
		case "-":
			// 区间
			if arrVal, ok := val.([]any); ok {
				if len(arrVal) < 2 {
					// paraName <= ?
					strVal, _ := arrVal[0].(string)
					strVal = strings.ReplaceAll(strVal, "-", "")
					ph := getPlaceholder(k, strVal)
					sqlParts = append(sqlParts, fmt.Sprintf("%s <= %s", paraName, ph))
					paramMap[ph] = strVal
				} else {
					// (paraName >= ? and paraName <= ?)
					strVal0, _ := arrVal[0].(string)
					strVal0 = strings.ReplaceAll(strVal0, "-", "")
					strVal1, _ := arrVal[1].(string)
					strVal1 = strings.ReplaceAll(strVal1, "-", "")

					ph0 := getPlaceholder(k, strVal0)
					ph1 := getPlaceholder(k, strVal1)
					sqlParts = append(sqlParts, fmt.Sprintf("(%s >= %s and %s <= %s)", paraName, ph0, paraName, ph1))
				}
			} else {
				strVal, _ := val.(string)
				if !strings.Contains(strVal, "-") {
					ph := getPlaceholder(k, strings.ReplaceAll(strVal, "-", ""))
					sqlParts = append(sqlParts, fmt.Sprintf("%s <= %s", paraName, ph))
				}
			}
		case "in":
			// in
			if arrVal, ok := val.([]any); ok {
				placeholders := []string{}
				for _, one := range arrVal {
					ph := getPlaceholder(k, one)
					placeholders = append(placeholders, ph)
				}
				sqlParts = append(sqlParts, fmt.Sprintf("%s in(%s)", paraName, strings.Join(placeholders, ",")))
			}
		case "not in":
			// not in
			if arrVal, ok := val.([]any); ok {
				placeholders := []string{}
				for _, one := range arrVal {
					ph := getPlaceholder(k, one)
					placeholders = append(placeholders, ph)
				}
				sqlParts = append(sqlParts, fmt.Sprintf("%s not in(%s)", paraName, strings.Join(placeholders, ",")))
			}
		case "is null":
			sqlParts = append(sqlParts, fmt.Sprintf("%s %s ?", paraName, spOp))
		case "is not null":
			sqlParts = append(sqlParts, fmt.Sprintf("%s %s ?", paraName, spOp))
		default:
			// no special op or other operators
			if spOp != "" {
				// paraName + spOps[k] + ?
				ph := getPlaceholder(k, val)
				sqlParts = append(sqlParts, fmt.Sprintf("%s %s %s", paraName, spOp, ph))
			} else {
				// in or default =
				if arrVal, ok := val.([]any); ok {
					placeholders := []string{}
					for _, one := range arrVal {
						ph := getPlaceholder(k, one)
						placeholders = append(placeholders, ph)
					}
					sqlParts = append(sqlParts, fmt.Sprintf("%s in(%s)", paraName, strings.Join(placeholders, ",")))
				} else {
					ph := getPlaceholder(k, val)
					sqlParts = append(sqlParts, fmt.Sprintf("%s = %s", paraName, ph))
				}
			}
		}
	}

	return strings.Join(sqlParts, " and "), paramMap
}

// func QueryParamsBuilder(params map[string]any, defaultPrefix string, otherPrefixes map[string][]string, spOps map[string]string) string {
// 	return NewQueryParamsBuilder().Params(params).DefaultPrefix(defaultPrefix).OtherPrefixes(otherPrefixes).SpOps(spOps).Build()
// }

func RegisterSqlHelper(dbType string, helper ISqlHelper) {
	g_DBHelpers[dbType] = helper
}
