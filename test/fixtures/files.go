package fixtures

import "madaurus/dev/material/app/models"

func GetFiles() []models.Files {

	files := []models.Files{}
	files = append(files, models.Files{
		Group: "Group 2",
		Name:  "Course OOP",
		Type:  "PDF",
	})
	files = append(files, models.Files{
		Group: "Group 1",
		Name:  "Course ACSI",
		Type:  "PDF",
	})
	files = append(files, models.Files{
		Group: "Group 1",
		Name:  "Course THL",
		Type:  "png",
	})
	files = append(files, models.Files{
		Group: "Group 1",
		Name:  "Course XML",
		Type:  "docs",
	})
	return files
}
