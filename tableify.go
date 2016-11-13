package tableify

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
)

type StructTable struct {
	*tablewriter.Table
	headerMap map[string]string
	headers   []string
	rows      [][]string
}

// API Design
//
// structTable := tableify.NewStructTable
// structTable.Append(data) --> where data is a struct
// structTable.AppendBulk(data)
// structTable.Render()
//
// StructTags
// Add StructTags to fields to control how they will be
// displayed via tableify
// eg type Example struct {
//		Value string `table:"-"`
// }
// Adding the "-" value to the table tag means ignore.
// You can also change the header name
// Or change that display values

func NewStructTable(writer io.Writer, typ interface{}) (*StructTable, error) {
	typType := reflect.TypeOf(typ)
	if typType.Kind() != reflect.Struct {
		return nil, errors.New("Given type is not a struct")
	}
	// typVal := reflect.ValueOf(typ)

	st := &StructTable{
		Table: tablewriter.NewWriter(writer),
	}

	var headers []string
	numFields := typType.NumField()
	for i := 0; i < numFields; i++ {
		field := typType.Field(i)
		header, _ := structFieldToLabel(field)
		if header != "" {
			headers = append(headers, header)
		}
	}

	// fmt.Println(headers)
	st.headers = headers
	st.SetHeader(headers)

	return st, nil
}

func structFieldToLabel(field reflect.StructField) (string, string) {
	tag := field.Tag.Get("table")
	var headerName string
	valueFormat := "{{.}}"
	if tag != "" {
		if strings.Contains(tag, ",") {
			parts := strings.Split(tag, ",")
			// fmt.Println(parts)
			if parts[0] == "-" {
				headerName = parts[0]
			} else if parts[0] == "" {
				headerName = field.Name
			} else {
				headerName = parts[0]
			}
			if parts[1] != "" {
				// fmt.Println("Header format:", headerName, parts[1])
				valueFormat = parts[1]
			}
		} else if tag == "-" {
			headerName = ""
		} else {
			headerName = tag
		}
	} else {
		headerName = field.Name
	}

	return headerName, valueFormat
}

func (t *StructTable) Append(data interface{}) {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Struct {
		return
	}

	dataVal := reflect.ValueOf(data)
	var values []string
	numFields := dataVal.NumField()
	for i := 0; i < numFields; i++ {
		headerName, valueFormat := structFieldToLabel(dataType.Field(i))
		if headerName != "" {
			value := dataVal.Field(i).Interface()
			valueStr := renderValue(valueFormat, value)
			values = append(values, valueStr)
		}
	}

	t.Table.Append(values)
	t.rows = append(t.rows, values)
}

func (t *StructTable) AppendBulk(data interface{}) {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice && dataType.Kind() != reflect.Array {
		return // silent failure. Need to fix
	}

	dataVal := reflect.ValueOf(data)

	length := dataVal.Len()
	for i := 0; i < length; i++ {
		elem := dataVal.Index(i)
		t.Append(elem.Interface())
	}
}

// take the format defined from the struct tag or default, and render the format
func renderValue(format string, value interface{}) string {
	output := bytes.NewBuffer([]byte{})
	// fmt.Println(format)
	t, _ := template.New("renderValue").Parse(format)
	t.Execute(output, value)
	str := output.String()
	// fmt.Println(str)
	return str
}
