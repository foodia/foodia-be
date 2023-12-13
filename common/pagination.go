package common

import (
	"gorm.io/gorm"
	"math"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
)

type Pagination struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	PageCount int   `json:"page_count"`
	Total     int64 `json:"total"`
	Next      int   `json:"next,omitempty"`
	Previous  int   `json:"previous,omitempty"`
}

func Paginate(query *gorm.DB, model any, p *Pagination) func(db *gorm.DB) *gorm.DB {
	var totalData int64
	query.Model(model).Count(&totalData)

	p.Total = totalData
	p.PageCount = int(math.Ceil(float64(totalData) / float64(p.PerPage)))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	return p.PerPage
}
