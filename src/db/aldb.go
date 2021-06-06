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
	success := mapResult(container, column, rs)
	if rows.Next() {
		panic("QueryOne except one result but get no more one")
	}
	return success
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
	//自定义sql 表达式中 ？由[]:变量名]代替，找到这些变量名并由反射根据改名称获取所给
	//结构体实例当中的数据作为参数传递给Exec函数
	reg, _ := regexp.Compile(`:[a-zA-z_]+`)
	regFind := reg.FindAllString(sqlStr, -1)
	//通过反射创建参数列表的容器
	params := make([]interface{}, len(regFind))
	//通过自定义sql表达式获取sql
	SQLParsed := reg.ReplaceAllString(sqlStr, "?")
	//通过自定义sql中：找到对应的参数
	for i, sqlArgs := range regFind {
		parseArg := strings.TrimPrefix(sqlArgs, `:`)
		fieldName := toPascalCase(parseArg)
		field := pointTo.FieldByName(fieldName)
		switch field.Kind() {
		case reflect.Int:
			//将参数添加到参数容器中
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

//mapResult 将sql rows扫描到的数据填入给定的结构中（结构体或slice)
//container :单条结果容器，columns 结果集对应数据库中的列名，value
//被映射对象
func mapResult(container []interface{}, columns []string, value reflect.Value) bool {
	var slot reflect.Value
	var arr = make([]reflect.Value, 0)
	//判断待映射类型，结构以与slice分别处理
	if value.Elem().Kind() == reflect.Struct {
		slot = value.Elem()
	} else {
		//slice内数据类型的实例
		slot = reflect.New(value.Type().Elem().Elem()).Elem()
	}
	var oneMoreSet = false
	//遍历一行结果集找到其在结构体中的位置并赋值
	for i, v := range container {
		//找到对应结构体的属性
		slotField := slot.FieldByName(toPascalCase(columns[i]))
		if slotField.CanSet() {
			switch value := v.(type) {
			case *int:
				//只有与其结构体类型匹配才赋值
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
	//如果被映射对象是slice也就是多结果集映射要通过反射将映射出的
	//结构体实例追加到结果集中
	if value.Elem().Kind() == reflect.Slice {
		arr = append(arr, slot)
		added := reflect.Append(value.Elem(), arr...)
		value.Elem().Set(added)
	}
	return oneMoreSet
}

//mapRes 将查询的结果集按一对多形式映射到结构当中
//allRows 所有结果集，columns 结果集对应数据库中的列名,value
//被映射对象,height工具属性与可变参数pk配合使用，pk（primary
//key）设计目的是为了兼容QueryOne与Query的结果集映射。实际
//这两个方法有单独的映射函数
func mapRes(allRows [][]interface{}, columns []string, value reflect.Value, height int, pk ...string) {
	in := value.Elem()
	inType := in.Type().Elem()
	var inSlot reflect.Value
	var inSlotName string
	//查找给定结构的slice属性并为其
	for i := 0; i < inType.NumField(); i++ {
		if inType.Field(i).Type.Kind() == reflect.Slice {
			//记录改属性属性名方便之后通过反射获取改属性并为其赋值
			inSlotName = inType.Field(i).Name
			inSlot = reflect.New(inType.Field(i).Type)
			mapRes(allRows, columns, inSlot, height+1, pk...)
		}
	}
	//mark为一个标识，以sql primary key为map，通过它标识同一元素是否被重复扫描
	mark := make(map[interface{}]byte)
	//主键在column中索引位置，方便获取主键值并配合mark判断是否重复扫描
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
}
func catchPanic() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}
