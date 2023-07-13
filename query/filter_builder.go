package query

type FilterBuilder struct {
	Where  []Filter
	Having []Filter
}

func (b *FilterBuilder) Build() (string, []interface{}) {
	// Create where query
	where, whereArgs := appendQuery(ClauseWhere, b.Where)
	// Append having query
	having, havingArgs := appendQuery(ClauseHaving, b.Having)
	// Merge query
	q := where + " " + having
	// Merge args
	a := append(whereArgs, havingArgs...)
	// Returns
	return q, a
}

func (b *FilterBuilder) Build2() (string, []interface{}) {
	// Create where query
	where, whereArgs := appendQuery(ClauseWhere2, b.Where)
	// Append having query
	having, havingArgs := appendQuery(ClauseHaving, b.Having)
	// Merge query
	q := where + " " + having
	// Merge args
	a := append(whereArgs, havingArgs...)
	// Returns
	return q, a
}

func appendQuery(clause string, filters []Filter) (string, []interface{}) {
	// Init args
	// Get filter length
	l := len(filters)
	// If filters length is empty, return empty string
	if l <= 0 {
		return "", nil
	}
	// Init query
	q := clause + " "
	// If filters length is exactly one, return first value
	if l == 1 {
		f := filters[0]
		return q + f.Query(), f.Args()
	}
	// Pop first value
	f, filters := filters[0], filters[1:]
	// Init query and args
	q += f.Query()
	a := f.Args()
	// Append other queries and args
	for _, f := range filters {
		q += " AND " + f.Query()
		a = append(a, f.Args()...)
	}
	// Return query
	return q, a
}

// NewFilterBuilder populates array of filter to each filter clause
func NewFilterBuilder(filters map[string]Filter) FilterBuilder {
	// init temporary slice
	var w []Filter
	var h []Filter
	// Extract map
	for _, f := range filters {
		c := f.QueryClause()
		switch c {
		case ClauseWhere:
			w = append(w, f)
		case ClauseHaving:
			h = append(h, f)
		}
	}
	// Return builder
	return FilterBuilder{w, h}
}
