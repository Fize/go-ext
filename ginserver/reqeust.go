package ginserver

const (
	_defaultPageSize    = 20
	_defaultCurrentPage = 1
)

const (
	// desc sort order
	Desc = "desc"
	// asc sort order
	Asc = "asc"
)

// Request request object
type Request struct {
	ID    uint64 `param:"id"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	// sort field, default is create time
	Sort string `form:"sort"`
	// sort order
	Order string `form:"order"`
}

func (q *Request) Default() {
	if q.Limit <= 0 {
		q.Limit = -1
		q.Page = 1
	}
}

// HandleQueryParam handle query params
func (q *Request) HandleQueryParam(total int) int {
	totalPages := 1
	if q.Limit < 0 {
		q.Page = totalPages
		q.Limit = total
		return totalPages
	}
	if q.Page <= 0 {
		q.Page = _defaultCurrentPage
	}
	if q.Limit <= 0 {
		q.Limit = _defaultPageSize
	}
	if total > q.Limit {
		totalPages = total / q.Limit
		if total%q.Limit > 0 {
			totalPages = totalPages + 1
		}
	}
	if q.Page > totalPages {
		q.Page = totalPages
	}

	if q.Order != Desc && q.Order != Asc {
		q.Order = Desc
	}
	return totalPages
}
