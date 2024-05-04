package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Section struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" validate:"required" binding:"required" bson:"name" `
	TeacherId string             `json:"teacher_id" bson:"teacher_id"`
	Videos    []Video            `json:"videos" bson:"videos"`
	Lectures  []Lecture          `json:"lectures" bson:"lectures"`
	Files     []Files            `json:"files" bson:"files"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}

type Lecture struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Group   string             `json:"group" bson:"group" binding:"required"`
	Name    string             `json:"name"`
	Content string             `json:"content" bson:"content"  binding:"required,min=250"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" validate:"required" binding:"required"`
	TeacherId string    `json:"teacher_id" bson:"teacher_id" validate:"required" binding:"required"`
	IsPublic  bool      `json:"is_public" bson:"is_public"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
type Video struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Group string             `json:"group" bson:"group" binding:"required"`

	Url       string             `json:"url" bson:"url" validate:"required"`
	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required"`
	TeacherId string             `json:"teacher_id" bson:"teacher_id" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}
type Files struct {
	ID  primitive.ObjectID `json:"id" bson:"_id"`
	Url string             `json:"url" bson:"url"`

	Name string `json:"name" bson:"name" validate:"required" binding:"required"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id"`
	Type      string    `json:"type" bson:"type"  binding:"default=text"`
	TeacherId string    `json:"teacher_id" bson:"teacher_id"`
	Group     string    `json:"group" bson:"group" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type StudentNote struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	StudentID string `json:"student_id" bson:"student_id" binding:"required"`
	//	SectionId primitive.ObjectID `json:"section_id" bson:"section_id" binding:"required" binding:"required"`
	Content string `json:"content" bson:"content" validate:"required" binding:"required"`
}
type ExtendedSection struct {
	Section

	Files    []Files        `json:"files"`
	Videos   []Video        `json:"videos"`
	Lectures []Lecture      `json:"contents"`
	Notes    *[]StudentNote `json:"note"`
}

// ---------------------- API ----------------------

type RSection struct {
	Section
	ID       string     `json:"id"`
	Files    []RFiles   `json:"files"`
	Videos   []RVideo   `json:"videos"`
	Lectures []RLecture `json:"lectures"`
}

func (s *RSection) Extract(section Section) {
	s = &RSection{
		Section:  section,
		ID:       section.ID.Hex(),
		Videos:   []RVideo{},
		Files:    []RFiles{},
		Lectures: []RLecture{},
	}
	for _, f := range section.Files {
		s.Files = append(s.Files, RFiles{
			f,
			f.ID.Hex(),
		})
	}
	for _, v := range section.Videos {
		s.Videos = append(s.Videos, RVideo{
			v,
			v.ID.Hex(),
		})
	}
	for _, l := range section.Lectures {
		s.Lectures = append(s.Lectures, RLecture{
			l,
			l.ID.Hex(),
		})
	}

}

// ---------- Files
type RFiles struct {
	Files
	ID string `json:"id"`
}

func (rf *RFiles) Extract(f Files) {
	rf = &RFiles{
		f,
		f.ID.Hex(),
	}

}

// ----------- Videos
type RVideo struct {
	Video
	ID string `json:"id"`
}

func (vd *RVideo) Extract(v Video) {
	vd = &RVideo{
		v,
		v.ID.Hex(),
	}
}

// --------- Lectures
type RLecture struct {
	Lecture
	ID string `json:"id"`
}

func (rl *RLecture) Extract(l Lecture) {
	rl = &RLecture{
		l,
		l.ID.Hex(),
	}
}
