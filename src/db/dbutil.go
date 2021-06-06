package db

import (
	"database/sql"
)

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
	case "VARCHAR", "TEXT", "LONGTEXT":
		return &sql.NullString{String: "", Valid: true}
	case "BIT":
		return new(interface{})
	case "DATE":
		return &sql.NullString{String: "", Valid: false}
	case "DATETIME":
		return &sql.NullString{String: "", Valid: false}
	case "TIMESTAMP":
		return &sql.NullString{String: "", Valid: false}
	}
	return nil
}

func toPascalCase(src string) string {
	var dst = make([]uint8, 0)
	if src[0] > 96 && src[0] < 123 {
		dst = append(dst, src[0]-32)
	} else {
		dst = append(dst, src[0])
	}
	for i := 1; i < len(src); {
		if src[i] == '_' {
			if src[0] > 96 && src[0] < 123 {
				dst = append(dst, src[i+1]-32)
			}
			i += 2
		} else {
			dst = append(dst, src[i])
			i++

		}
	}
	return string(dst)
}
func getColIndex(colunms []string, col string) int {
	for idx, item := range colunms {
		if item == col {
			return idx
		}
	}
	return -1
}
func pkValue(pkContent interface{}) interface{} {
	switch v := pkContent.(type) {
	case *int:
		return *v
	case *byte:
		return *v
	case *float32:
		return *v
	case *string:
		return *v
	case *sql.NullString:
		return v.String
	default:
		return nil
	}
}
