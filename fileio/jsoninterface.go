package fileio

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

const (
	OBJ_PARENT_ATTR_NAME  = "PARENT"
	PARENT_ATTR_NOT_FOUND = "___PARENT_ATTR_NOT_FOUND___"
	ATTR_NO_MEMBER_OBJ    = "_n|m_"
)

func JsonToSpreadsheet(fp string, parentAttr string) ([][]string, error) {
	rec := [][]string{}
	header := []string{}

	if parentAttr != "" {
		header = append(header, OBJ_PARENT_ATTR_NAME)
	}

	content, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	if !json.Valid(content) {
		return nil, fmt.Errorf("json file is not valid")
	}

	// var res map[string]interface{}
	var res interface{}
	json.Unmarshal(content, &res)

	var objs []map[string]string
	switch reflect.ValueOf(res).Kind().String() {
	case "slice":
		parseJsonSlice(res.([]interface{}), &objs, &header, parentAttr, "")
	case "map":
		parseJsonMap(res.(map[string]interface{}), &objs, &header, parentAttr, "")
	default:
		return nil, fmt.Errorf("malformed json!? Expecting array or map!")
	}

	fmt.Println("Found Objects: ", len(objs))
	for _, obj := range objs {
		fmt.Printf("%s (%s)\n", obj["_diProId"], obj[OBJ_PARENT_ATTR_NAME])
	}
	rec = mapSliceToSpreadsheet(objs, header)
	fmt.Println(len(rec))
	for _, row := range rec {
		fmt.Println(row[0], row[1])
	}
	return rec, nil
}

func parseJsonSlice(l []interface{}, objs *[]map[string]string, header *[]string, parentAttr string, parent string) {
	for _, v := range l {
		switch reflect.ValueOf(v).Kind().String() {
		case "slice":
			parseJsonSlice(v.([]interface{}), objs, header, parentAttr, parent)
		case "map":
			parseJsonMap(v.(map[string]interface{}), objs, header, parentAttr, parent)
		default:
			fmt.Println("--->>> ", reflect.ValueOf(v).Kind().String())
		}
	}
}

func parseJsonMap(m map[string]interface{}, objs *[]map[string]string, header *[]string, parentAttr string, parent string) {
	obj := make(map[string]string)
	curParVal := ""
	if parentAttr != "" {
		parVal, ok := m[parentAttr]
		if !ok {
			fmt.Println("Parent-Attribute not found")
		} else {
			curParVal = parVal.(string)
		}
	}
	obj[OBJ_PARENT_ATTR_NAME] = parent
	for k, v := range m {
		switch reflect.ValueOf(v).Kind().String() {
		case "slice":
			parseJsonSlice(v.([]interface{}), objs, header, parentAttr, curParVal)
		case "map":
			parseJsonMap(v.(map[string]interface{}), objs, header, parentAttr, curParVal)
		default:
			if hasValue(*header, k) < 0 {
				*header = append(*header, k)
			}
			obj[k] = fmt.Sprintf("%v", v)
		}
	}
	*objs = append(*objs, obj)
}

func hasValue(l []string, v string) int {
	for i, e := range l {
		if e == v {
			return i
		}
	}
	return -1
}

func mapSliceToSpreadsheet(objs []map[string]string, header []string) [][]string {
	rec := [][]string{header}
	for k := len(objs) - 1; k >= 0; k-- {
		row := make([]string, len(header))
		for i, attr := range header {
			v, ok := objs[k][attr]
			if ok {
				row[i] = v
			} else {
				row[i] = ATTR_NO_MEMBER_OBJ
			}
		}
		rec = append(rec, row)
	}
	return rec
}
