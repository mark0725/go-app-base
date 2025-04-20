package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/godror/godror"
	applog "github.com/mark0725/go-app-base/logger"
	"github.com/mark0725/go-app-base/utils"
)

type DatabaseConnection struct {
	Name   string
	Type   string
	Config *DatabaseConfig
	Conn   *sql.DB
}

type IDatabase interface {
	Init(db string, config *DatabaseConfig) error
	//DBExec(sql string, params []any) (any, error)
	//DBQuery(sql string, params []any) ([]map[string]any, error)
	GetDBConn(name string) *DatabaseConnection
	GetDBTables(db string, name string) ([]map[string]any, error)
	GetDBTableFields(db string, table string) ([]map[string]any, error)
}

var logger = applog.GetLogger("database")
var g_DBConns map[string]*DatabaseConnection = make(map[string]*DatabaseConnection)

const DB_CONN_NAME_DEFAULT string = "default"

func InitDB(name string, config *DatabaseConfig) error {
	// 配置 MySQL 数据库连接参数
	logger = applog.GetLogger("database")

	dsn := ""
	dbType := ""
	driverName := ""
	if config.Type == "" {
		dbType = "mysql"
		driverName = "mysql"
	} else {
		dbType = config.Type
	}

	if dbType == "mysql" && config.Driver == "" {
		driverName = "mysql"
	} else {
		driverName = config.Driver
	}

	switch dbType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.DBUser,
			config.DBPass,
			config.Host,
			config.Port,
			config.DBName,
		)
		if config.Options != "" {
			dsn += "?" + config.Options
		}
	case "oracle":
		dsn = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`,
			config.DBUser,
			config.DBPass,
			config.Host,
			config.Port,
			config.DBName,
		)
		if config.Options != "" {
			dsn += " " + config.Options
		}
	case "postgres":
		if config.Driver == "" {
			driverName = "pgx"
		}
		//connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			config.DBUser,
			config.DBPass,
			config.Host,
			config.Port,
			config.DBName,
		)

		if config.Options != "" {
			dsn += "?" + config.Options
		} else {
			dsn += "?sslmode=disable"
		}

	default:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBPass,
			config.Host,
			config.Port,
			config.DBName,
		)
	}

	logger.Tracef("dsn: %s", dsn)
	// 打开数据库连接
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		logger.Error(err)
		return err
	}
	//defer db.Close()

	// 验证连接是否成功
	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return err
	}

	conn := &DatabaseConnection{
		Name:   name,
		Type:   dbType,
		Config: config,
		Conn:   db,
	}

	logger.Info("database connect success.")
	if name != "" {
		g_DBConns[name] = conn
	} else {
		g_DBConns[DB_CONN_NAME_DEFAULT] = conn
	}
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(0)

	// fmt.Println("Successfully connected to MySQL database")

	// // 执行查询
	// rows, err := db.Query("SELECT id, name FROM example_table")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// // 迭代查询结果
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	if err := rows.Scan(&id, &name); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("id: %d, name: %s\n", id, name)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	return nil
}

func DBExec(db string, sql string, params []any) (any, error) {
	conn, exist := g_DBConns[db]
	if !exist {
		logger.Errorf("not found db: %s", db)
		return nil, errors.New("not found db: " + db)
	}

	result, err := conn.Conn.Exec(sql, params...)
	if err != nil {
		logger.Errorf("Sql Error: %v: %s, %v", err, sql, params)
		return nil, err
	}

	// insertedID, _ := result.LastInsertId()
	logger.Tracef("Sql success: %s, %v", sql, params)

	return result, nil
}

func DBPageQuery(db string, sql string, params map[string]any, pageSize int, pageIndex int) (*PageQueryResult, error) {

	result := &PageQueryResult{
		Content:          []map[string]any{},
		First:            true,
		Last:             true,
		Number:           0,
		NumberOfElements: 0,
		Size:             pageSize,
		TotalElements:    0,
		TotalPages:       0,
	}

	if pageIndex > 0 {
		result.First = false
	}

	countSql := fmt.Sprintf("select count(*) num from (%s) _PAGECOUNT", sql)
	sqlCounterBuildResult := QueryNamedParamsBuilder(countSql, params)
	logger.Trace("count query sql:", sqlCounterBuildResult.Sql)
	logger.Trace("count query params:", sqlCounterBuildResult.Params)
	dbCountResult, err := DBQuery(db, sqlCounterBuildResult.Sql, sqlCounterBuildResult.Params)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	if len(dbCountResult) == 0 {
		result.Last = true
		return result, nil
	}

	num := dbCountResult[0]["num"].(int64)
	if num == 0 {
		result.Last = true
		return result, nil
	}

	logger.Trace("row count:", num)
	result.TotalElements = int(num)
	result.TotalPages = int(math.Ceil(float64(result.TotalElements) / float64(result.Size)))
	logger.Trace("page count:", result.TotalPages)
	if result.TotalPages == 0 || pageIndex == result.TotalPages-1 {
		result.Last = true
	}

	startIndex := pageIndex * pageSize

	querySql, _ := PageQueryBuilder(db, sql, pageSize, startIndex)
	sqlBuildResult := QueryNamedParamsBuilder(querySql, params)
	logger.Trace("query sql:", sqlBuildResult.Sql)
	logger.Trace("query params:", sqlBuildResult.Params)
	dbResult, err := DBQuery(db, sqlBuildResult.Sql, sqlBuildResult.Params)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	if len(dbResult) == 0 {
		result.Last = true
		return result, nil
	}

	result.Content = dbResult

	result.Number = pageIndex
	result.NumberOfElements = len(result.Content)

	return result, nil
}

func DBQuery(db string, sql string, params []any) ([]map[string]any, error) {
	conn, exist := g_DBConns[db]
	if !exist {
		logger.Errorf("not found db: %s", db)
		return nil, errors.New("not found db: " + db)
	}

	// 执行查询
	rows, err := conn.Conn.Query(sql, params...)
	if err != nil {
		logger.Errorf("Sql Error: %v: %s, %v", err, sql, params)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		logger.Errorf("Sql Result Col Error: %v: %s, %v", err, sql, params)
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	var records []map[string]any

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logger.Errorf("Sql Result Row Error: %v: %s, %v", err, sql, params)
			return nil, err
		}

		rowData := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}

		records = append(records, rowData)
	}

	return records, nil

}

func GetDBConn(name string) *DatabaseConnection {
	connName := DB_CONN_NAME_DEFAULT
	if name != "" {
		connName = name
	}

	if conn, exists := g_DBConns[connName]; exists {
		return conn
	}

	return nil
}

func GetDBTables(db string, name string) ([]map[string]any, error) {
	conn := GetDBConn(db)
	if conn == nil {
		return nil, errors.New("not found db: " + db)
	}
	switch conn.Type {
	case "mysql":
		return MysqlGetDBTables(db, name)
	case "postgres":
		return PostgresGetDBTables(db, name)
	case "oracle":
		return OracleGetDBTables(db, name)
	case "sqlite":
		return SqliteGetDBTables(db, name)
	}

	return nil, errors.New("not found db: " + db)
}

func GetDBTableFields(db string, table string) ([]map[string]any, error) {
	conn := GetDBConn(db)
	if conn == nil {
		return nil, errors.New("not found db: " + db)
	}

	switch conn.Type {
	case "mysql":
		return MysqlGetDBTableFields(db, table)
	case "postgres":
		return PostgresGetDBTableFields(db, table)
	case "oracle":
		return OracleGetDBTableFields(db, table)
	case "sqlite":
		return SqliteGetDBTableFields(db, table)
	}

	return nil, errors.New("not found db: " + db)
}

func DBType2GoType(db string, dbType string) (string, error) {
	conn := GetDBConn(db)
	if conn == nil {
		return "", errors.New("not found db: " + db)
	}

	switch conn.Type {
	case "mysql":
		return MysqlDBType2GoType(dbType)
	case "postgres":
		return PostgresDBType2GoType(dbType)
	case "oracle":
		return OracleDBType2GoType(dbType)
	case "sqlite":
		return SqliteDBType2GoType(dbType)
	default:
		return "", errors.New("not found db type: " + conn.Type)
	}

}

func DBQueryEnt[T any](db string, table string, where string, params map[string]any) ([]*T, error) {
	var ent T
	fields, err := GetTableFieldIds(&ent)
	if err != nil {
		logger.Error("GetTableFieldIds fail: ", err)
		return nil, err
	}

	whereSql := ""
	if where != "" {
		whereSql = " where " + where
	}

	sqlOld := "SELECT " + strings.Join(fields, ",") + " FROM  " + table + whereSql

	sqlResult := QueryNamedParamsBuilder(sqlOld, params)
	logger.Trace("query sql:", sqlResult.Sql)
	logger.Trace("query params:", sqlResult.Params)

	result, err := DBQuery(db, sqlResult.Sql, sqlResult.Params)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	if len(result) == 0 {
		logger.Trace("DBQuery notfound: ")
		return nil, nil
	}

	rows := utils.Map(result, func(row map[string]any) *T {
		var rec T
		if err := MapRowToStruct(row, &rec); err != nil {
			logger.Error("MapRowToStruct fail: ", err)
			return nil
		}

		return &rec
	})

	return rows, nil
}

func DBInsertEnt[T any](db string, table string, ent *T) error {
	// fields, err := GetTableFieldIds(ent)
	// if err != nil {
	// 	logger.Error("GetTableFieldIds fail: ", err)
	// 	return err
	// }

	params := EntiryToMap(ent)
	sqlhelper := GetSqlHelper(db)
	if sqlhelper == nil {
		return errors.New("not found sqlhelper: " + db)
	}

	insertSql := sqlhelper.InsertSqlBuilder(table, params, nil, nil)
	_, err := DBExec(db, insertSql.Sql, insertSql.Params)

	return err
}
