package db

import (
	"fmt"
	"strings"
	//_ "github.com/go-sql-driver/mysql"
)

var mysqlToGoTypeMap = map[string]string{
	"tinyint":    "int8",
	"smallint":   "int16",
	"mediumint":  "int32",
	"int":        "int64",
	"integer":    "int",
	"bigint":     "int64",
	"float":      "float32",
	"double":     "float64",
	"decimal":    "float64",
	"char":       "string",
	"varchar":    "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"binary":     "[]byte",
	"varbinary":  "[]byte",
	"blob":       "[]byte",
	"tinyblob":   "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"date":       "time.Time",
	"time":       "time.Time",
}

func MysqlInitDB(name string, config *DatabaseConfig) error {
	// TODO: Implement MySQL database initialization
	return nil
}

func MysqlGetDBTables(db string, filter string) ([]map[string]any, error) {
	filterStr := ""
	if filter != "" {
		filterStr = fmt.Sprintf(" and table_name like '%s'", filter)
	}

	result, err := DBQuery(db, `SELECT table_schema, table_name, table_type, engine FROM information_schema.tables WHERE 1=1`+filterStr, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func MysqlGetDBTableFields(db string, tableName string) ([]map[string]any, error) {
	logger.Debug("genEntity:", tableName)
	sql := fmt.Sprintf(`SELECT column_name, data_type, column_type, column_default, is_nullable, column_key, extra, column_comment FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION ASC`, tableName)
	result, err := DBQuery(db, sql, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func MysqlDBType2GoType(DBType string) (string, error) {
	// 去掉可能存在的长度和其他限定符，如 int(11)
	typeParts := strings.Split(DBType, "(")
	cleanedType := typeParts[0]

	goType, found := mysqlToGoTypeMap[strings.ToLower(cleanedType)]
	if !found {
		return "", fmt.Errorf("unsupported DB type: %s", DBType)
	}

	return goType, nil
}
