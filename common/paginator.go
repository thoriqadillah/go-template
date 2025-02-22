package common

import (
	"math"
	"net/url"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type Pager struct {
	Page      int   `json:"page"`
	TotalPage int   `json:"totalPage"`
	Limit     int   `json:"limit"`
	Count     int   `json:"count"`
	Items     []any `json:"items"`
}

type Paginator struct {
	params url.Values
	Page   int
	Limit  int
	Offset int
}

func Paginate(params url.Values) *Paginator {
	page := Parse(params.Get("page")).Int(1)
	limit := Parse(params.Get("limit")).Int(10)
	offset := (page - 1) * limit

	return &Paginator{
		params: params,
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func (p *Paginator) totalPage(count int) int {
	total := math.Ceil(float64(count / p.Limit))
	if int(total) <= 0 {
		return 1
	}

	return int(total)
}

func (p *Paginator) CreatePager(items []any, count int) Pager {
	return Pager{
		Page:      p.Page,
		Items:     items,
		Count:     count,
		TotalPage: p.totalPage(count),
		Limit:     p.Limit,
	}
}

type FilterFunc func(v Parser) bob.Mod[*dialect.SelectQuery]

type paramBuilder struct {
	paginator Paginator
	mods      []bob.Mod[*dialect.SelectQuery]
}

func FilterParam(p Paginator) *paramBuilder {
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
