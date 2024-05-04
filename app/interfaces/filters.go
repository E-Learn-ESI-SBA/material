package interfaces

type ModuleFilter struct {
	Year       *int8   `json:"year,omitempty"`
	Semester   *int8   `json:"semester,omitempty"`
	Speciality *string `json:"speciality,omitempty"`
}

type PaginationQuery struct {
	Page  int `json:"page,omitempty"`
	Items int `json:"items" `
}

func (p *PaginationQuery) newPagination(page int, items int) {
	p.Page = page
	p.Items = items
}

func (m *ModuleFilter) newModuleFilter(year int8, semester int8, speciality string) {
	m.Year = &year
	m.Semester = &semester
	m.Speciality = &speciality
}
