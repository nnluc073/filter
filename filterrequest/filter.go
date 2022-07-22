package filterrequest

import "fmt"

type FilterReq struct {
	Data map[string]interface{}
}

// Has check if the given field exists in the request data.
func (r *FilterReq) Has(field string) bool {
	_, exists := r.Data[field]
	return exists
}

// String get a string field from the request data.
// Panics if the field is not a string.
func (r *FilterReq) String(field string) string {
	str, ok := r.Data[field].(string)
	if !ok {
		panic(fmt.Sprintf("Field \"%s\" is not a string", field))
	}
	return str
}

// Numeric get a numeric field from the request data.
// Panics if the field is not numeric.
func (r *FilterReq) Numeric(field string) float64 {
	str, ok := r.Data[field].(float64)
	if !ok {
		panic(fmt.Sprintf("Field \"%s\" is not numeric", field))
	}
	return str
}

// Integer get an integer field from the request data.
// Panics if the field is not an integer.
func (r *FilterReq) Integer(field string) int {
	str, ok := r.Data[field].(int)
	if !ok {
		panic(fmt.Sprintf("Field \"%s\" is not an integer", field))
	}
	return str
}

// Bool get a bool field from the request data.
// Panics if the field is not a bool.
func (r *FilterReq) Bool(field string) bool {
	str, ok := r.Data[field].(bool)
	if !ok {
		panic(fmt.Sprintf("Field \"%s\" is not a bool", field))
	}
	return str
}

// Object get an object field from the request data.
// Panics if the field is not an object.
func (r *FilterReq) Object(field string) map[string]interface{} {
	str, ok := r.Data[field].(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Field \"%s\" is not an object", field))
	}
	return str
}
