package common

import (
	"math"
	"net/url"
	"strings"

	"github.com/uptrace/bun"
)

type Pager struct {
	Page      int   `json:"page"`
	TotalPage int   `json:"totalPage"`
	Limit     int   `json:"limit"`
	Count     int64 `json:"count"`
	Items     any   `json:"items"`
}

type Paginator struct {
	params *url.Values
	page   int
	Limit  int
	Offset int
}

func Paginate(params url.Values) *Paginator {
	page := Parse(params.Get("page")).Int(1)
	limit := Parse(params.Get("limit")).Int(5)
	offset := (page - 1) * limit

	return &Paginator{
		params: &params,
		page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func (p *Paginator) totalPage(count int64) int {
	total := math.Ceil(float64(count / int64(p.Limit)))
	if int(total) <= 0 {
		return 1
	}

	return int(total)
}

func (p *Paginator) CreatePager(items any, count int64) Pager {
	return Pager{
		Page:      p.page,
		Items:     items,
		Count:     count,
		TotalPage: p.totalPage(count),
		Limit:     p.Limit,
	}
}

type FilterFunc func(v Parser, sq *bun.SelectQuery) *bun.SelectQuery

type SearchFunc func(search string, sq *bun.SelectQuery) *bun.SelectQuery

type paramBuilder struct {
	paginator *Paginator
	query     *bun.SelectQuery
}

func FilterParam(p *Paginator, baseQuery *bun.SelectQuery) *paramBuilder {
	return &paramBuilder{
		paginator: p,
		query:     baseQuery,
	}
}

func (p *paramBuilder) Search(filter SearchFunc) *paramBuilder {
	key := "search"
	if !p.paginator.params.Has(key) {
		return p
	}

	param := p.paginator.params.Get(key)
	p.query.Apply(func(sq *bun.SelectQuery) *bun.SelectQuery {
		return filter(param, sq)
	})

	return p
}

func (p *paramBuilder) SearchOn(columns ...string) *paramBuilder {
	key := "search"
	if !p.paginator.params.Has(key) {
		return p
	}

	param := p.paginator.params.Get(key)

	for i, column := range columns {
		p.query.Apply(func(sq *bun.SelectQuery) *bun.SelectQuery {
			if i == 0 {
				return sq.Where("? ILIKE ?", bun.Ident(column), "%"+param+"%")
			}

			return sq.WhereOr("? ILIKE ?", bun.Ident(column), "%"+param+"%")
		})
	}

	return p
}

func (p *paramBuilder) ApplyOrder() {
	key := "orderBy"
	if !p.paginator.params.Has(key) {
		return
	}

	order := strings.Split(p.paginator.params.Get(key), ":")
	p.query.Apply(func(sq *bun.SelectQuery) *bun.SelectQuery {
		if strings.ToLower(order[1]) == "desc" {
			return sq.OrderExpr("? DESC", bun.Ident(order[0]))
		}

		return sq.OrderExpr("? ASC", bun.Ident(order[0]))
	})
}

// Filter will run the filter function when there is param key in the url
func (p *paramBuilder) Filter(key string, filter FilterFunc) *paramBuilder {
	if !p.paginator.params.Has(key) {
		return p
	}

	param := p.paginator.params.Get(key)
	p.query.Apply(func(sq *bun.SelectQuery) *bun.SelectQuery {
		return filter(Parse(param), sq)
	})

	return p
}
