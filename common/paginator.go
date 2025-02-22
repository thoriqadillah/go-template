package common

import (
	"math"
	"net/url"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type pager struct {
	Page      int   `json:"page"`
	TotalPage int   `json:"totalPage"`
	Limit     int   `json:"limit"`
	Count     int   `json:"count"`
	Items     []any `json:"items"`
}

type paginator struct {
	params url.Values
	page   int
	limit  int
	offset int
}

func Paginate(params url.Values) *paginator {
	page := Parse(params.Get("page")).Int(1)
	limit := Parse(params.Get("limit")).Int(10)
	offset := (page - 1) * limit

	return &paginator{
		params: params,
		page:   page,
		limit:  limit,
		offset: offset,
	}
}

func (p *paginator) totalPage(count int) int {
	total := math.Ceil(float64(count / p.limit))
	if int(total) <= 0 {
		return 1
	}

	return int(total)
}

func (p *paginator) CreatePaginator(slice []any, count int) pager {
	return pager{
		Page:      p.page,
		Items:     slice,
		Count:     count,
		TotalPage: p.totalPage(count),
		Limit:     p.limit,
	}
}

type FilterFunc func(v Parser) bob.Mod[*dialect.SelectQuery]

type paramBuilder struct {
	paginator *paginator
	mods      []bob.Mod[*dialect.SelectQuery]
}

func FilterParam(p *paginator) *paramBuilder {
	return &paramBuilder{
		paginator: p,
		mods:      make([]bob.Mod[*dialect.SelectQuery], 0),
	}
}

// Filter will run the filter function when there is param key in the url
func (p *paramBuilder) Filter(key string, filter FilterFunc) *paramBuilder {
	if !p.paginator.params.Has(key) {
		return p
	}

	param := p.paginator.params.Get(key)
	p.mods = append(p.mods, filter(Parse(param)))
	return p
}

func (p *paramBuilder) Build() []bob.Mod[*dialect.SelectQuery] {
	return p.mods
}
