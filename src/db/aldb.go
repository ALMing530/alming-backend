package db

import (
	"database/sql"
	"log"
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
		log.Println("An error occerred when exec query sql", err)
	}
	rc, err := rows.ColumnTypes()
	if err != nil {
		log.Println("Get column types fail")
	}
	container := createContainer(rc)
	if !rows.Next() {
		return false
	}
	column, _ := rows.Columns()
	rows.Scan(container...)
	return mapResult(container, column, rs)
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
		log.Println("An error occerred when exec query sql", err)
	}

	rc, err := rows.ColumnTypes()
	if err != nil {
		log.Println("Get column types fail")
	}

	column, _ := rows.Columns()
	var oneMoreSet bool = false
	for rows.Next() {
		container := createContainer(rc)
		err = rows.Scan(container...)
		if err != nil {
			panic("Scan rows error")
		}
		oneMoreSet = mapResult(container, column, rs)
	}
	return oneMoreSet
}
func QueryOneToMany(slice interface{}, sqlStr string, outPk string, inPk string, params ...interface{}) (resMatched bool) {
	defer catchPanic()
	rs := reflect.ValueOf(slice)
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
		log.Println("An error occerred when exec query sql", err)
	}

	rc, err := rows.ColumnTypes()
	if err != nil {
		log.Println("Get column types fail")
	}

	column, _ := rows.Columns()
	var allRows = make([][]interface{}, 0)
	for rows.Next() {
		container := createContainer(rc)
		err = rows.Scan(container...)
		if err != nil {
			panic("Scan rows error")
		}
		allRows = append(allRows, container)
	}
	mapRes(allRows, column, rs, 0, outPk, inPk)
	//别忘改
	return true
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
		fieldName := toPascalCase(parseArg)
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
	var res sql.Result
	var err error
	if len(params) > 0 {
		res, err = DB.Exec(SQLParsed, params...)
	} else {
		res, err = DB.Exec(SQLParsed)
	}
	if err == nil {
		rowAf, _ := res.RowsAffected()
		return rowAf > 0
	}
	return false
}
func mapResult(container []interface{}, columns []string, value reflect.Value) bool {
	var slot reflect.Value
	var arr = make([]reflect.Value, 0)
	if value.Elem().Kind() == reflect.Struct {
		slot = value.Elem()
	} else {
		slot = reflect.New(value.Type().Elem().Elem()).Elem()
	}
	var oneMoreSet = false
	for i, v := range container {
		slotField := slot.FieldByName(toPascalCase(columns[i]))
		if slotField.CanSet() {
			switch value := v.(type) {
			case *int:
				if slotField.Kind() == reflect.Int {
					slotField.SetInt(int64(*value))
				}
				oneMoreSet = true
			case *string:
				if slotField.Kind() == reflect.String {
					slotField.SetString(*value)
				}
				oneMoreSet = true
			case *sql.NullString:
				if slotField.Kind() == reflect.String {
					slotField.SetString(value.String)
				}
				oneMoreSet = true
			}
		}
	}
	if value.Elem().Kind() == reflect.Slice {
		arr = append(arr, slot)
		added := reflect.Append(value.Elem(), arr...)
		value.Elem().Set(added)
	}
	return oneMoreSet
}
func mapRes(allRows [][]interface{}, columns []string, value reflect.Value, height int, pk ...string) reflect.Value {
	in := value.Elem()
	inType := in.Type().Elem()
	var inSlot reflect.Value
	var inSlotName string
	for i := 0; i < inType.NumField(); i++ {
		if inType.Field(i).Type.Kind() == reflect.Slice {
			inSlotName = inType.Field(i).Name
			inSlot = reflect.New(inType.Field(i).Type)
			mapRes(allRows, columns, inSlot, height+1, pk...)
		}
	}
	mark := make(map[interface{}]byte)
	var pkIdx = -1
	if len(pk) > 0 {
		pkIdx = getColIndex(columns, pk[height])
	}
	var arr = make([]reflect.Value, 0)
	for _, row := range allRows {
		if mark[pkValue(row[pkIdx])] == 1 {
			continue
		}
		outSlot := reflect.New(inType).Elem()
		var oneMoreSet = false
		for i, v := range row {
			slot := outSlot.FieldByName(toPascalCase(columns[i]))
			if slot.CanSet() {
				switch setValue := v.(type) {
				case *int:
					if slot.Kind() == reflect.Int {
						slot.SetInt(int64(*setValue))
						oneMoreSet = true
					}
				case *string:
					if slot.Kind() == reflect.String {
						slot.SetString(*setValue)
						oneMoreSet = true
					}
				case *sql.NullString:
					if slot.Kind() == reflect.String {
						slot.SetString(setValue.String)
						oneMoreSet = true
					}
				}
			}
		}
		slot := outSlot.FieldByName(inSlotName)
		if slot.CanSet() {
			slot.Set(inSlot.Elem())
		}
		if oneMoreSet {
			if len(pk) > 0 {
				mark[pkValue(row[pkIdx])] = 1
			}
		}
		arr = append(arr, outSlot)
	}
	added := reflect.Append(in, arr...)
	in.Set(added)
	return in
}
func catchPanic() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}
