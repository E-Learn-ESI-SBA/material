package fixtures

import "madaurus/dev/material/app/models"

func GetLectures() []models.Lecture {
	linkMark := models.Mark{
		Type: "link",
		Attrs: map[string]string{
			"href":   "https://example.com",
			"target": "_blank",
		},
	}

	// Create sample Nodes
	textNode1 := models.Node{
		Type: "text",
		Text: "This is a link",
		Marks: []models.Mark{
			linkMark,
		},
	}

	paragraphNode := models.Node{
		Type: "paragraph",
		Content: []models.Node{
			textNode1,
		},
	}

	headingNode := models.Node{
		Type:  "heading",
		Attrs: map[string]interface{}{"level": 2},
		Content: []models.Node{
			{Type: "text", Text: "Sample Heading"},
		},
	}

	// Create Content struct
	content := models.Content{
		Type: "doc",
		Content: []models.Node{
			headingNode,
			paragraphNode,
		},
	}
	var lectures []models.Lecture
	lectures = append(lectures, models.Lecture{
		Name:      "Lecture 1",
		Content:   content,
		TeacherId: "12",
		IsPublic:  false,
		Groups:    "Group 1",
	})
	lectures = append(lectures, models.Lecture{
		Name:      "Lecture 2",
		Groups:    "Group 1",
		Content:   content,
		IsPublic:  false,
		TeacherId: "12",
	})
	lectures = append(lectures, models.Lecture{
		Name:      "Lecture 2",
		Content:   content,
		IsPublic:  false,
		TeacherId: "15",
		Groups:    "Group 1",
	})
	return lectures
}
