package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//QueryOne Query data using given sql  and map to struct given
//The structure is struct that you want to map,the params is sql parmas
func QueryOne(structure interface{}, sqlStr string, params ...interface{}) (resMatched bool) {
	defer catchPanic()
	rs := reflect.ValueOf(structure)
	pointTo := rs.Elem()
	if pointTo.Kind() != reflect.Struct {
		panic("QueryOne must to map to a struct,please check your structure parameter")
	}
	var rows *sql.Rows
	var err error
	if len(params) == 0 {
		rows, err = DB.Query(sqlStr)
	} else {
		rows, err = DB.Query(sqlStr, params...)
	}
	if err != nil {
		fmt.Println("An error occerred when exec query sql", err)
	}
	rc, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println("Get column types fail")
	}
	container := createContainer(rc)
	if !rows.Next() {
		return false
	}
	column, _ := rows.Columns()
	rows.Scan(container...)
	if rows.Next() {
		panic("QueryOne accept one result but get no more one")
	}
	var oneMoreSet bool = false
	for i, v := range container {
		rField := pointTo.FieldByName(strings.Title(column[i]))
		if rField.CanSet() {
			switch v.(type) {
			case *int:
				if rField.Kind() == reflect.Int {
					rField.SetInt(int64(*v.(*int)))
				}
				oneMoreSet = true
			case *string:
				if rField.Kind() == reflect.String {
					rField.SetString(*v.(*string))
				}
				oneMoreSet = true
			}
		}
	}
	if !oneMoreSet {
		return false
	}
	return true
}

//Query Query data using given sql  and map to slice given
//The structure is slice that you want to map,the params is sql parmas
func Query(structure interface{}, sqlStr string, params ...interface{}) (resMatched bool) {
	defer catchPanic()
	rs := reflect.ValueOf(structure)
	pointTo := rs.Elem()
	if pointTo.Kind() != reflect.Slice {
		panic("QueryOne must to map to a slice,please check your structure parameter")
	}
	var rows *sql.Rows
	var err error
	if len(params) == 0 {
		rows, err = DB.Query(sqlStr)
	} else {
		rows, err = DB.Query(sqlStr, params...)
	}
	if err != nil {
		fmt.Println("An error occerred when exec query sql", err)
	}

	rc, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println("Get column types fail")
	}

	column, _ := rows.Columns()
	inType := rs.Type().Elem().Elem()
	var oneMoreSet bool = false
	temp := make([]reflect.Value, 0)
	for rows.Next() {
		container := createContainer(rc)
		rows.Scan(container...)
		slot := reflect.New(inType).Elem()
		for i, v := range container {
			rField := slot.FieldByName(strings.Title(column[i]))
			if rField.CanSet() {
				switch v.(type) {
				case *int:
					if rField.Kind() == reflect.Int {
						rField.SetInt(int64(*v.(*int)))
					}
					oneMoreSet = true
				case *string:
					if rField.Kind() == reflect.String {
						rField.SetString(*v.(*string))
					}
					oneMoreSet = true
				}
			}
		}
		temp = append(temp, slot)
	}
	arrAdded := reflect.Append(pointTo, temp...)
	pointTo.Set(arrAdded)
	return oneMoreSet
}

//Exec excute sql with the params in the struct you give
func Exec(structure interface{}, sqlStr string) (success bool) {
	rs := reflect.ValueOf(structure)
	pointTo := rs.Elem()
	reg, _ := regexp.Compile(`:[a-zA-z_]+`)
	regFind := reg.FindAllString(sqlStr, -1)
	params := make([]interface{}, len(regFind))
	SQLParsed := reg.ReplaceAllString(sqlStr, "?")
	for i, sqlArgs := range regFind {
		parseArg := strings.TrimPrefix(sqlArgs, `:`)
		fieldName := strings.Title(parseArg)
		field := pointTo.FieldByName(fieldName)
		switch field.Kind() {
		case reflect.Int:
			params[i] = field.Int()
		case reflect.String:
			params[i] = field.String()
		case reflect.Float32, reflect.Float64:
			params[i] = field.Float()
		}
	}
	if res, err := DB.Exec(SQLParsed, params...); err != nil {
		rowAf, _ := res.RowsAffected()
		return rowAf > 0
	}
	return false
}
func createContainer(columnTyes []*sql.ColumnType) (params []interface{}) {
	params = make([]interface{}, len(columnTyes))
	for i, ct := range columnTyes {
		params[i] = createSlot(ct.DatabaseTypeName())
	}
	return
}
func createSlot(dbType string) interface{} {
	switch dbType {
	case "INT", "TINYINT", "BIGINT":
		return new(int)
	case "MEDIUMINT":
	case "DOUBLE":
		return new(float32)
	case "DECIMAL":
	case "CHAR":
		return new(byte)
	case "VARCHAR":
		return new(string)
	case "BIT":
		return new(interface{})
	case "DATE":
		return new(string)
	case "DATETIME":
		return new(string)
	case "TIMESTAMP":
		return new(string)
	}
	return nil
}
func catchPanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}