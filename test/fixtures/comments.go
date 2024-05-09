package fixtures

import "madaurus/dev/material/app/models"

func GetComment() []models.Comments {

	var comments []models.Comments
	comments = append(comments, models.Comments{
		Content:  "Hello world",
		IsEdited: false,
	})
	comments = append(comments, models.Comments{
		Content:  "Hello algeria",
		IsEdited: true,
	})
	return comments
}
