package db

import (
	"fmt"
	"strings"
	//_ "github.com/godror/godror" // Oracle driver
)

var oracleToGoTypeMap = map[string]string{
	"number":                         "float64", // Can be int or float based on scale, using float64 as safe default
	"integer":                        "int64",
	"float":                          "float64",
	"binary_float":                   "float32",
	"binary_double":                  "float64",
	"char":                           "string",
	"varchar":                        "string",
	"varchar2":                       "string",
	"nvarchar2":                      "string",
	"clob":                           "string",
	"blob":                           "[]byte",
	"raw":                            "[]byte",
	"date":                           "time.Time",
	"timestamp":                      "time.Time",
	"timestamp with time zone":       "time.Time",
	"timestamp with local time zone": "time.Time",
}

func OracleInitDB(name string, config *DatabaseConfig) error {
	// TODO: Implement Oracle database initialization
	return nil
}

func OracleGetDBTables(db string, filter string) ([]map[string]any, error) {
	filterStr := ""
	if filter != "" {
		filterStr = fmt.Sprintf(" AND table_name LIKE '%s'", strings.ToUpper(filter))
	}

	result, err := DBQuery(db, `SELECT owner AS table_schema, table_name, 'BASE TABLE' AS table_type, 
                              'N/A' AS engine FROM ALL_TABLES WHERE 1=1`+filterStr, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func OracleGetDBTableFields(db string, tableName string) ([]map[string]any, error) {
	logger.Debug("genEntity:", tableName)
	fields := []string{
		"column_name",
		"data_type",
		"data_type || '(' || data_length || ')' AS column_type",
		"data_default AS column_default",
		"nullable AS is_nullable",
		"decode(constraint_type, 'P', 'PRI', '') AS column_key",
		"'' AS extra", // Oracle doesn't have an exact equivalent to MySQL's EXTRA
		"comments AS column_comment",
	}

	// Join with ALL_COL_COMMENTS for comments and ALL_CONS_COLUMNS for constraints
	sql := fmt.Sprintf(`
		SELECT %s 
		FROM ALL_TAB_COLUMNS atc
		LEFT JOIN ALL_COL_COMMENTS acc 
			ON atc.owner = acc.owner 
			AND atc.table_name = acc.table_name 
			AND atc.column_name = acc.column_name
		LEFT JOIN (
			SELECT acc2.column_name, ac.constraint_type 
			FROM ALL_CONS_COLUMNS acc2
			JOIN ALL_CONSTRAINTS ac 
				ON acc2.constraint_name = ac.constraint_name 
				AND ac.constraint_type = 'P'
		) cons 
			ON atc.column_name = cons.column_name
		WHERE atc.table_name = UPPER('%s') 
		ORDER BY atc.column_id ASC`,
		strings.Join(fields, ","),
		tableName,
	)

	result, err := DBQuery(db, sql, nil)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	return result, nil
}

func OracleDBType2GoType(DBType string) (string, error) {
	// Remove possible length and other qualifiers
	typeParts := strings.Split(DBType, "(")
	cleanedType := typeParts[0]

	goType, found := oracleToGoTypeMap[strings.ToLower(cleanedType)]
	if !found {
		return "", fmt.Errorf("unsupported DB type: %s", DBType)
	}

	return goType, nil
}
