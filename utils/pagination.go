package utils

import "math"

type Pagination struct {
	CurrentPage int
	PerPage int
	Total int
}

func (p *Pagination) AllPages() (pages int){
	if p.PerPage == 0 {
		pages = 0
	} else {
		pages = int(math.Ceil(float64(p.Total )/ float64(p.PerPage) ))
	}
	return pages
}

func (p *Pagination) HasPrev() bool {
	return p.CurrentPage > 1
}

func (p *Pagination) PrevNum() int {
	if !p.HasPrev() {
		return -1
	}
	return p.CurrentPage - 1
}

func (p *Pagination) HasNext() bool{
	return p.CurrentPage < p.AllPages()
}

func (p *Pagination) NextNum() int {
	if !p.HasNext() {
		return -1
	}
	return p.CurrentPage + 1
}

func (p *Pagination) PageRet()[]int {
	var leftEdge int = 2
	var leftCurrent int = 2
	var rightCurrent int = 2
	var rightEdge int = 2
	var res []int
	last := 0
	allPages := p.AllPages()
	for i:=1;i<allPages+1;i++ {
		if i <= leftEdge || (i>p.CurrentPage - leftCurrent - 1 && i<p.CurrentPage + rightCurrent) || i> allPages - rightEdge {
			if last + 1!= i {
				res = append(res, -1)
			}
			res = append(res, i)
			last = i
		}
	}
	return res
}












