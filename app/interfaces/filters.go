package interfaces

type ModuleFilter struct {
	Year       *int8   `json:"year,omitempty" validate:"min=1,max=5"`
	Semester   *int8   `json:"semester,omitempty" validate:"min=1,max=2"`
	Speciality *string `json:"speciality,omitempty"`
}

type PaginationQuery struct {
	Page  int8 `json:"page,omitempty" validate:"min=1 ,default=1"`
	Items int8 `json:"items" validate:"min=1, default=10"`
}

func (p *PaginationQuery) newPagination(page int8, items int8) {
	p.Page = page
	p.Items = items
}

func (m *ModuleFilter) newModuleFilter(year int8, semester int8, speciality string) {
	m.Year = &year
	m.Semester = &semester
	m.Speciality = &speciality
}
