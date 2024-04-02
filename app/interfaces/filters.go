package interfaces

type ModuleFilter struct {
	Year       *int8   `json:"year,omitempty" validate:"min=1,max=5"`
	Semester   *int8   `json:"semester,omitempty" validate:"min=1,max=2"`
	Speciality *string `json:"speciality,omitempty"`
}

type PaginationQuery struct {
	Page  int `json:"page,omitempty" validate:"min=1 ,default=1"`
	Items int `json:"items" validate:"min=1, default=10"`
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
