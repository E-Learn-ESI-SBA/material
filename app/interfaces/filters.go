package interfaces

type CourseFilter struct {
	Year       int8   `json:"year" validate:"min=1,max=5"`
	Semester   int8   `json:"semester" validate:"min=1,max=2"`
	Speciality string `json:"speciality,omitempty"`
}
