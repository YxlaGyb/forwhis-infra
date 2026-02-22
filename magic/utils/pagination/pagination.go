package pagination

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type PageRequest struct {
	Page  int
	Limit int
}

type PageResult struct {
	Total  int64
	Page   int
	Limit  int
	Offset int
	Pages  int
}

func NewPageRequest(page, limit int) *PageRequest {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return &PageRequest{
		Page:  page,
		Limit: limit,
	}
}

func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

func NewPageResult(total int64, page, limit int) *PageResult {
	if limit < 1 {
		limit = DefaultLimit
	}
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}
	return &PageResult{
		Total:  total,
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
		Pages:  pages,
	}
}

func NormalizePage(page int) int {
	if page < 1 {
		return DefaultPage
	}
	return page
}

func NormalizeLimit(limit int) int {
	if limit < 1 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

func CalcOffset(page, limit int) int {
	return (NormalizePage(page) - 1) * NormalizeLimit(limit)
}
