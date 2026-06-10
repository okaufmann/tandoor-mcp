package tandoor

import (
	"fmt"
	"net/url"
)

// QueryBuilder is a simple wrapper around url.Values to make appending query parameters easy.
type QueryBuilder struct {
	v url.Values
}

// NewQuery returns a new QueryBuilder
func NewQuery() *QueryBuilder {
	return &QueryBuilder{v: url.Values{}}
}

// Add appends a query parameter to the builder if it is not a zero/nil value.
func (q *QueryBuilder) Add(key string, value any) *QueryBuilder {
	switch v := value.(type) {
	case string:
		if v != "" {
			q.v.Add(key, v)
		}
	case *string:
		if v != nil && *v != "" {
			q.v.Add(key, *v)
		}
	case int:
		// Assuming 0 is a valid int to send or we want to send it. 
		// If 0 should be omitted, the caller should use *int instead.
		q.v.Add(key, fmt.Sprintf("%d", v))
	case *int:
		if v != nil {
			q.v.Add(key, fmt.Sprintf("%d", *v))
		}
	case []int:
		for _, val := range v {
			q.v.Add(key, fmt.Sprintf("%d", val))
		}
	case []string:
		for _, val := range v {
			if val != "" {
				q.v.Add(key, val)
			}
		}
	case bool:
		q.v.Add(key, fmt.Sprintf("%t", v))
	case *bool:
		if v != nil {
			q.v.Add(key, fmt.Sprintf("%t", *v))
		}
	default:
		if v != nil {
			q.v.Add(key, fmt.Sprintf("%v", v))
		}
	}
	return q
}

// Values returns the underlying url.Values
func (q *QueryBuilder) Values() url.Values {
	return q.v
}
