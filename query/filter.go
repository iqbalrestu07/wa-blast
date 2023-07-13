package query

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"wa-blast/util"
)

const (
	ClauseWhere  = ""
	ClauseWhere2 = ""
	ClauseHaving = "HAVING"

	FilterBool               = 1
	FilterFloat32            = 2
	FilterInt32              = 3
	FilterInt8               = 8
	FilterRadius             = 4
	FilterSince              = 5
	FilterStringLike         = 6
	FilterStringMatch        = 7
	FilterStringLikeMultiple = 9
)

// Regexp
var whitespaceRegexp = regexp.MustCompile(" +")

// Filter is an interface that get filter from http request query and converts it to SQL query
type Filter interface {
	Args() []interface{}
	Query() string
	QueryClause() string
	Set(string) error
}

// sqlString is string extension to implement Filter setter
type sqlString struct {
	Clause string
	String string
}

func (q *sqlString) Query() string {
	return q.String
}

func (q *sqlString) QueryClause() string {
	return q.Clause
}

// RadiusFilter is an implementation of query.Filter for distance radius filter
type RadiusFilter struct {
	Radius int32
	sqlString
}

func (f *RadiusFilter) Set(value string) error {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	// If radius <= 0, return error
	if val <= 0 {
		return errors.New("radius must be > 0")
	}
	// Set radius
	f.Radius = int32(val)
	return nil
}

func (f *RadiusFilter) Args() []interface{} {
	return []interface{}{f.Radius}
}

// StringMatchFilter is an implementation of query.Filter for string matching filter
type StringMatchFilter struct {
	Value string
	sqlString
}

func (f *StringMatchFilter) Set(value string) error {
	f.Value = value
	return nil
}

func (f *StringMatchFilter) Args() []interface{} {
	return []interface{}{f.Value}
}

// StringLikeFilter is an implementation of query.Filter for string like filter
type StringLikeFilter struct {
	Value string
	sqlString
}

func (f *StringLikeFilter) Set(value string) error {
	f.Value = value
	return nil
}

func (f *StringLikeFilter) Args() []interface{} {
	// Convert white space to %
	val := whitespaceRegexp.ReplaceAllString(f.Value, "%")
	return []interface{}{"%" + val + "%"}
}

// StringLikeFilter is an implementation of query.Filter for string like filter
type StringLikeMultipleFilter struct {
	Value string
	sqlString
}

func (f *StringLikeMultipleFilter) Set(value string) error {
	f.Value = value
	return nil
}

func (f *StringLikeMultipleFilter) Args() []interface{} {
	// Convert white space to %
	val := whitespaceRegexp.ReplaceAllString(f.Value, "%")
	return []interface{}{"%\"" + val + "\"%"}
}

// Float32Filter is an implementation of query.Filter for float based filter
type Float32Filter struct {
	Value float32
	sqlString
}

func (f *Float32Filter) Set(value string) error {
	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return err
	}
	// Set value
	f.Value = float32(val)
	// Return
	return nil
}

func (f *Float32Filter) Args() []interface{} {
	return []interface{}{f.Value}
}

// Int32Filter is an implementation of query.Filter for float based filter
type Int32Filter struct {
	Value int32
	sqlString
}

func (f *Int32Filter) Set(value string) error {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	// Set value
	f.Value = int32(val)
	// Return
	return nil
}

func (f *Int32Filter) Args() []interface{} {
	return []interface{}{f.Value}
}

// Int8Filter is an implementation of query.Filter for float based filter
type Int8Filter struct {
	Value int8
	sqlString
}

func (f *Int8Filter) Set(value string) error {
	val, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return err
	}
	// Set value
	f.Value = int8(val)
	// Return
	return nil
}

func (f *Int8Filter) Args() []interface{} {
	return []interface{}{f.Value}
}

// BoolFilter is an implementation of query.Filter for boolean data type filter
type BoolFilter struct {
	Value bool
	sqlString
}

func (f *BoolFilter) Set(value string) error {
	if value == "true" || value == "1" {
		f.Value = true
	} else {
		f.Value = false
	}
	return nil
}

func (f *BoolFilter) Args() []interface{} {
	return []interface{}{f.Value}
}

// SinceFilter is an implementation of query.Filter for since filter
type SinceFilter struct {
	Value int64
	sqlString
}

func (f *SinceFilter) Set(value string) error {
	f.Value = util.ParseInt64(value, 0)
	// Return
	return nil
}

func (f *SinceFilter) Args() []interface{} {
	return []interface{}{f.Value}
}

// NewFilter construct filter by its type
func NewFilter(filterType int8, clause string, query string) Filter {
	// init sqlString
	s := sqlString{clause, query}
	var f Filter
	switch filterType {
	case FilterBool:
		f = &BoolFilter{sqlString: s}
	case FilterFloat32:
		f = &Float32Filter{sqlString: s}
	case FilterInt8:
		f = &Int8Filter{sqlString: s}
	case FilterInt32:
		f = &Int32Filter{sqlString: s}
	case FilterRadius:
		f = &RadiusFilter{sqlString: s}
	case FilterSince:
		f = &SinceFilter{sqlString: s}
	case FilterStringLike:
		f = &StringLikeFilter{sqlString: s}
	case FilterStringMatch:
		f = &StringMatchFilter{sqlString: s}
	case FilterStringLikeMultiple:
		f = &StringLikeMultipleFilter{sqlString: s}
	default:
		fmt.Println("unknown filter type")
		os.Exit(24)
	}
	return f
}
