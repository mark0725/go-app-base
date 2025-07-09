package db

import (
	"fmt"
	"strings"
)

// PostgreSQL type mapping to Go types, including pgvector support
var pgToGoTypeMap = map[string]string{
	"smallint":         "int32",
	"integer":          "int64",
	"bigint":           "int64",
	"real":             "float32",
	"double precision": "float64",
	"numeric":          "float64",
	"decimal":          "float64",
	"character":        "string",
	"varchar":          "string",
	"char":             "string",
	"bpchar":           "string",
	"text":             "string",
	"bytea":            "[]byte",
	"timestamp":        "time.Time",
	"timestamptz":      "time.Time",
	"date":             "time.Time",
	"time":             "time.Time",
	"timetz":           "time.Time",
	"boolean":          "bool",
	"vector":           "[]float32", // pgvector 的 vector 类型映射到 Go 的 float32 切片
	"halfvec":          "[]float32", // pgvector 的 vector 类型映射到 Go 的 float32 切片
}

func PostgresInitDB(name string, config *DatabaseConfig) error {
	// 初始化 PostgreSQL 数据库并启用 pgvector 扩展
	_, err := DBQuery(name, "CREATE EXTENSION IF NOT EXISTS vector", nil)
	if err != nil {
		logger.Error("Failed to enable pgvector extension: ", err)
		return err
	}
	return nil
}

func PostgresGetDBTables(db string, filter string) ([]map[string]any, error) {
	filterStr := ""
	if filter != "" {
		filterStr = fmt.Sprintf(" AND table_name ILIKE '%%%s%%'", filter)
	}

	result, err := DBQuery(db, `
		SELECT table_schema, 
		       table_name, 
		       table_type
		FROM information_schema.tables 
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		`+filterStr, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func PostgresGetDBTableFields(db string, tableName string) ([]map[string]any, error) {
	logger.Debug("genEntity:", tableName)
	sql := `SELECT
			c.column_name,
			CASE 
				WHEN c.udt_name = 'bpchar' THEN 'char'
				WHEN c.udt_name = 'varchar' THEN 'varchar'
				WHEN c.udt_name = 'int4' THEN 'integer'
				WHEN c.udt_name = 'int8' THEN 'bigint'
				WHEN c.udt_name = 'numeric' THEN 'numeric'
				WHEN c.udt_name = 'text' THEN 'text'
				WHEN c.udt_name = 'bool' THEN 'boolean'
				WHEN c.udt_name = 'timestamp' THEN 'timestamp'
				WHEN c.udt_name = 'date' THEN 'date'
				WHEN c.udt_name = 'time' THEN 'time'
				ELSE c.udt_name
			END AS data_type,
			CASE 
				WHEN c.data_type = 'character varying' THEN 'VARCHAR(' || c.character_maximum_length || ')'
				WHEN c.data_type = 'numeric' THEN 'NUMERIC(' || c.numeric_precision || ',' || c.numeric_scale || ')'
				ELSE 
					CASE 
						WHEN c.udt_name = 'bpchar' THEN 'CHAR(' || c.character_maximum_length || ')'
						WHEN c.udt_name = 'varchar' THEN 'VARCHAR(' || c.character_maximum_length || ')'
						WHEN c.udt_name = 'int4' THEN 'INTEGER'
						WHEN c.udt_name = 'int8' THEN 'BIGINT'
						WHEN c.udt_name = 'numeric' THEN 'NUMERIC(' || c.numeric_precision || ',' || c.numeric_scale || ')'
						WHEN c.udt_name = 'text' THEN 'TEXT'
						WHEN c.udt_name = 'bool' THEN 'BOOLEAN'
						WHEN c.udt_name = 'timestamp' THEN 'TIMESTAMP'
						WHEN c.udt_name = 'date' THEN 'DATE'
						WHEN c.udt_name = 'time' THEN 'TIME'
						ELSE c.udt_name
					END
			END AS column_type,
			c.column_default,
			c.is_nullable,
			(CASE WHEN kcu.column_name IS NOT NULL THEN 'PRI' END) AS column_key,
			(CASE 
				WHEN c.column_default LIKE 'nextval(%::regclass)' THEN 'auto_increment'
				ELSE ''
			END) AS extra,
			pd.description AS column_comment
		FROM
			information_schema.columns c
		LEFT JOIN
			pg_catalog.pg_class cl ON cl.relname = c.table_name
		LEFT JOIN
			pg_catalog.pg_namespace n ON n.nspname = c.table_schema AND n.oid = cl.relnamespace
		LEFT JOIN
			pg_catalog.pg_attribute a ON a.attrelid = cl.oid AND a.attname = c.column_name
		LEFT JOIN
			pg_catalog.pg_description pd ON pd.objoid = cl.oid AND pd.objsubid = a.attnum
		LEFT JOIN
			information_schema.key_column_usage kcu 
			ON c.table_schema = kcu.table_schema 
			AND c.table_name = kcu.table_name 
			AND c.column_name = kcu.column_name
			AND kcu.constraint_name IN (
				SELECT constraint_name 
				FROM information_schema.table_constraints 
				WHERE constraint_type = 'PRIMARY KEY'
			)
		WHERE
			c.table_schema = 'public' 
			AND c.table_name = $1
		ORDER BY
			c.ordinal_position;
	`

	result, err := DBQuery(db, sql, []interface{}{tableName})
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func PostgresDBType2GoType(DBType string) (string, error) {
	// 处理可能的数组类型和限定符
	cleanedType := strings.TrimSuffix(strings.ToLower(DBType), "[]")

	// 检查是否为 vector 类型
	if strings.HasPrefix(cleanedType, "vector") {
		return "[]float32", nil // pgvector 的 vector 类型映射为 []float32
	}

	goType, found := pgToGoTypeMap[cleanedType]
	if !found {
		// 处理一些常见的模糊匹配情况
		switch {
		case strings.Contains(cleanedType, "character"):
			return "string", nil
		case strings.Contains(cleanedType, "numeric"):
			return "float64", nil
		case strings.Contains(cleanedType, "timestamp"):
			return "time.Time", nil
		default:
			return "", fmt.Errorf("unsupported PostgreSQL type: %s", DBType)
		}
	}

	return goType, nil
}
