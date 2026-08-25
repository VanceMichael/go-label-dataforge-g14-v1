package pagination

type Query struct {
	Limit, Offset int
	Sort, Filter  string
}

func Normalize(q Query) Query {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	return q
}
