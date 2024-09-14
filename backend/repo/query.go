package repo

type Query struct {
	Limit   int
	Offset  int
	OrderBy string
	Asc     bool
}

func NewQuery(limit int, offset int, orderBy string, Asc bool) Query {
	if limit == 0 {
		limit = 10
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	return Query{
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy,
		Asc:     Asc,
	}
}
