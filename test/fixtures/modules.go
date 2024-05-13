package fixtures

import (
	"madaurus/dev/material/app/models"
)

func GetModules() []models.Module {
	var modules []models.Module
	image := "https://images.unsplash.com/photo-1617854818583-09e7f077a156?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
	// generate dummy data
	speciality := "Siw"
	speciality2 := "ISI"
	instructors := []string{"1", "2"}

	modules = append(modules, models.Module{
		Name:        "Module 1",
		Year:        "1",
		Speciality:  &speciality,
		Semester:    1,
		Coefficient: 1,
		TeacherId:   "1",
		Instructors: &instructors,
		IsPublic:    true,
		Plan:        []string{"Plan 1", "Plan 2"},
		Image:       &image,
	})
	modules = append(modules, models.Module{
		Name:        "Module 2",
		Year:        "2",
		Speciality:  &speciality2,
		IsPublic:    true,
		Semester:    2,
		Coefficient: 2,
		TeacherId:   "2",
		Instructors: &instructors,
		Plan:        []string{"Plan 1", "Plan 2"},
		Image:       &image,
	})
	modules = append(modules, models.Module{
		Name:        "Module 3",
		Year:        "3",
		Speciality:  &speciality,
		Semester:    2,
		IsPublic:    false,
		Image:       &image,
		Coefficient: 3,
		Plan:        []string{"Plan 1", "Plan 2"},
		Instructors: &instructors,
		TeacherId:   "53",
	})
	return modules
}
