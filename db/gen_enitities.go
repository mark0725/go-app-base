package db

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"unicode"
)

func GenEntities(db string, filter string, dir string, maps map[string]string) error {
	m := maps
	if m == nil {
		m = make(map[string]string)
	}

	result, err := GetDBTables(db, filter)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return err
	}

	for _, row := range result {
		tableName := row["table_name"].(string)
		entityName := ToCamelCase(tableName)
		if name, ok := m[tableName]; ok {
			entityName = ToCamelCase(name)
		}

		err := genEntity(db, tableName, entityName, dir)
		if err != nil {
			logger.Error("genEntity fail: ", err)
			return err
		}
	}

	return nil
}

type EntityField struct {
	Name          string `field-id:"column_name"`
	DataType      string `field-id:"data_type"`
	IsNullable    string `field-id:"is_nullable"`
	ColumnKey     string `field-id:"column_key"`
	ColumnDefault string `field-id:"column_default"`
	ColumnType    string `field-id:"column_type"`
	Extra         string `field-id:"extra"`
	Comment       string `field-id:"column_comment"`
}

func genEntity(db string, tableName string, entityName string, dir string) error {
	logger.Debug("genEntity:", tableName)
	result, err := GetDBTableFields(db, tableName)
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return err
	}

	tableFields := make([]EntityField, 0)

	for _, row := range result {
		rowObj := EntityField{}
		MapRowToStruct(row, &rowObj)
		tableFields = append(tableFields, rowObj)

		logger.Trace("entity:", rowObj)
	}

	enitityContent := "package entities \n\n"
	enitityContent += "\n"
	enitityContent += fmt.Sprintf(`const DB_TABLE_%s = "%s"`, strings.ToUpper(tableName), tableName) + "\n\n"
	enitityContent += fmt.Sprintf("type %s struct {", entityName)

	for _, field := range tableFields {
		varName := ToCamelCase(field.Name)
		varType, err := DBType2GoType(db, field.DataType)
		if err != nil {
			logger.Error("DBType2GoType fail: ", err)
			return err
		}

		fieldComment := field.Comment
		if fieldComment != "" {
			fieldComment = fmt.Sprintf(" // %s", fieldComment)
		}
		enitityContent += fmt.Sprintf("\n%s %s"+"`"+`field-id:"%s" field-type:"%s"`+"`%s", varName, varType, field.Name, field.DataType, fieldComment)
	}
	enitityContent += "\n}"

	//将内容写入文件
	filePath := path.Join(dir, strings.ToLower(tableName)+".go")
	out, err := os.Create(filePath)
	if err != nil {
		logger.Error("Create file fail: ", err)
		return err
	}
	defer out.Close()

	io.WriteString(out, enitityContent)

	return nil

}

func ToCamelCase(input string) string {
	// 将输入字符串全转为小写
	input = strings.ToLower(input)

	// 使用 strings.Builder 更高效的字符串拼接
	var result strings.Builder
	// 将字符串分割为子字符串
	words := strings.Split(input, "_")

	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		// 将第一个字母转换为大写
		for i, r := range word {
			if i == 0 {
				result.WriteRune(unicode.ToUpper(r))
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}
