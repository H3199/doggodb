package data

import (
	"fmt"
)
// A row is a slice of values
type Row struct {
    Values []interface{}
}

func CreateRow(values []interface{}) Row {
    return Row{Values: values}
}

func (r *Row) GetValue(index int) (interface{}, error) {
    if index < 0 || index >= len(r.Values) {
        return nil, fmt.Errorf("index out of bounds")
    }
    return r.Values[index], nil
}

func (r *Row) SetValue(index int, value interface{}) error {
    if index < 0 || index >= len(r.Values) {
        return fmt.Errorf("index out of bounds")
    }
    r.Values[index] = value
    return nil
}
