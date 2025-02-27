package db

import (
	"fmt"
	"strings"
	//_ "github.com/mattn/go-sqlite3"
)

var sqliteToGoTypeMap = map[string]string{
	"integer":  "int64",
	"int":      "int64",
	"smallint": "int16",
	"tinyint":  "int8",
	"bigint":   "int64",
	"real":     "float64",
	"double":   "float64",
	"float":    "float32",
	"numeric":  "float64",
	"text":     "string",
	"varchar":  "string",
	"char":     "string",
	"blob":     "[]byte",
	"datetime": "time.Time",
	"date":     "time.Time",
	"time":     "time.Time",
}

func SqliteInitDB(name string, config *DatabaseConfig) error {
	// TODO: Implement SQLite database initialization
	return nil
}

func SqliteGetDBTables(db string, filter string) ([]map[string]any, error) {
	filterStr := ""
	if filter != "" {
		filterStr = fmt.Sprintf(" AND name LIKE '%s'", filter)
	}

	// SQLite 使用 sqlite_master 表来查询表信息
	result, err := DBQuery(db, `SELECT name AS table_name, type AS table_type, sql FROM sqlite_master WHERE type='table'`+filterStr, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func SqliteGetDBTableFields(db string, tableName string) ([]map[string]any, error) {
	logger.Debug("genEntity:", tableName)

	// SQLite 使用 PRAGMA table_info 获取列信息
	sql := fmt.Sprintf(`PRAGMA table_info('%s')`, tableName)
	result, err := DBQuery(db, sql, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	// SQLite 的 PRAGMA table_info 返回的字段包括：
	// cid, name, type, notnull, dflt_value, pk
	// 我们需要将这些转换为与 MySQL 版本相似的格式
	for _, row := range result {
		// 重命名字段以保持兼容性
		if name, ok := row["name"]; ok {
			row["column_name"] = name
			delete(row, "name")
		}
		if typ, ok := row["type"]; ok {
			row["data_type"] = typ
			delete(row, "type")
		}
		if dflt, ok := row["dflt_value"]; ok {
			row["column_default"] = dflt
			delete(row, "dflt_value")
		}
		if nn, ok := row["notnull"]; ok {
			isNullable := "YES"
			if nn.(int64) == 1 {
				isNullable = "NO"
			}
			row["is_nullable"] = isNullable
			delete(row, "notnull")
		}
		if pk, ok := row["pk"]; ok {
			pkVal := pk.(int64)
			if pkVal == 1 {
				row["column_key"] = "PRI"
			} else {
				row["column_key"] = ""
			}
			delete(row, "pk")
		}
		// SQLite 没有 comment 字段，添加空值以保持兼容性
		row["column_comment"] = ""
		// 添加 column_type 以匹配 MySQL 版本
		row["column_type"] = row["data_type"]
		// extra 字段在 SQLite 中留空
		row["extra"] = ""
	}

	return result, nil
}

func SqliteDBType2GoType(DBType string) (string, error) {
	// 去掉可能存在的长度和其他限定符
	typeParts := strings.Split(DBType, "(")
	cleanedType := typeParts[0]

	goType, found := sqliteToGoTypeMap[strings.ToLower(cleanedType)]
	if !found {
		// SQLite 是动态类型系统，许多类型可以映射到基本类型
		switch strings.ToLower(cleanedType) {
		case "boolean":
			return "bool", nil
		case "character", "nchar", "nvarchar":
			return "string", nil
		case "decimal":
			return "float64", nil
		default:
			return "", fmt.Errorf("unsupported DB type: %s", DBType)
		}
	}

	return goType, nil
}
