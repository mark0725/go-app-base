package db

import (
	"database/sql"
	"errors"
	"fmt"

	applog "github.com/mark0725/go-app-base/logger"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger = nil
var g_DBConns map[string]*sql.DB = make(map[string]*sql.DB)

const DB_CONN_NAME_DEFAULT string = "default"

func InitDB(name string, config *DatabaseConfig) error {
	// 配置 MySQL 数据库连接参数
	logger = applog.GetLogger("database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.DBUser,
		config.DBPass,
		config.Host,
		config.Port,
		config.DBName,
	)

	logger.Tracef("dsn: %s", dsn)
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
		return err
	}
	//defer db.Close()

	// 验证连接是否成功
	err = db.Ping()
	if err != nil {
		log.Error(err)
		return err
	}

	logger.Info("database connect success.")
	if name != "" {
		g_DBConns[name] = db
	} else {
		g_DBConns[DB_CONN_NAME_DEFAULT] = db
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

func DBExec(db string, sql string, params []interface{}) (interface{}, error) {
	conn, exist := g_DBConns[db]
	if !exist {
		logger.Errorf("not found db: %s", db)
		return nil, errors.New("not found db: " + db)
	}

	result, err := conn.Exec(sql, params...)
	if err != nil {
		logger.Errorf("Sql Error: %v: %s, %v", err, sql, params)
		return nil, err
	}

	// insertedID, _ := result.LastInsertId()
	logger.Tracef("Sql success: %s, %v", sql, params)

	return result, nil
}

func DBQuery(db string, sql string, params []interface{}) ([]map[string]interface{}, error) {
	conn, exist := g_DBConns[db]
	if !exist {
		logger.Errorf("not found db: %s", db)
		return nil, errors.New("not found db: " + db)
	}

	// 执行查询
	rows, err := conn.Query(sql, params...)
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

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	var records []map[string]interface{}

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logger.Errorf("Sql Result Row Error: %v: %s, %v", err, sql, params)
			return nil, err
		}

		rowData := make(map[string]interface{})
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

func GetDBConn(name string) *sql.DB {
	connName := DB_CONN_NAME_DEFAULT
	if name != "" {
		connName = name
	}

	if conn, exists := g_DBConns[connName]; exists {
		return conn
	}

	return nil
}
