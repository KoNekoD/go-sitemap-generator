package stm

import "fmt"

// Attrs defines for xml attribute.
type Attrs []any

// Attr defines for xml attribute.
type Attr map[string]string

// URL User should use this typedef in main func.
type URL [][]any

func (u URL) URLJoinBy(key string, joins ...string) URL {
	var values []string
	for _, k := range joins {
		var vals any
		for _, v := range u {
			if v[0] == k {
				vals = v[1]
				break
			}
		}
		values = append(values, fmt.Sprint(vals))
	}
	var index int
	var v []any
	for index, v = range u {
		if v[0] == key {
			break
		}
	}
	u[index][1] = URLJoin("", values...)
	return u
}
