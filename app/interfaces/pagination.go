package interfaces

type PaginationQuery struct {
	page  int8
	items int8
}

func (p *PaginationQuery) newPagination(page int8, items int8) {
	p.page = page
	p.items = items
}
