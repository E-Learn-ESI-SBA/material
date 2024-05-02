package fixtures

import "madaurus/dev/material/app/models"

func GetSections() []models.Section {
	var sections []models.Section
	sections = append(sections, models.Section{
		Name: "Section 1",
	},
		models.Section{
			Name: "Section 2",
		},
		models.Section{
			Name: "Section 3",
		},
		models.Section{
			Name: "Section 4",
		},
	)

	return sections

}
