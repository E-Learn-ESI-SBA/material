package quizes

import (
	"errors"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Create Quiz
// @Description Protected Route used to create a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param quiz body models.Quiz true "Quiz Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Failure 500 {object} interfaces.APiError
// @Router /quizes/create [POST]
func CreateQuiz(collection *mongo.Collection, moduleCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var quiz models.Quiz
		err := c.BindJSON(&quiz)
		user := c.MustGet("user").(*utils.UserDetails)
		quiz.TeacherId = user.ID
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.CreateQuiz(c.Request.Context(), collection, moduleCollection, quiz)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Quiz Created Successfully"})
	}
}



// @Summary Update Quiz
// @Description Protected Route used to update a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param quiz body models.Quiz true "Quiz Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/update [PUT]
func UpdateQuiz(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var quiz models.Quiz
		user := c.MustGet("user").(*utils.UserDetails)
		err := c.BindJSON(&quiz)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		updatedAt := time.Now()
		quiz.UpdatedAt = &updatedAt
		err = services.UpdateQuiz(c.Request.Context(), collection, quiz, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Quiz Updated Successfully"})
	}

}



// @Summary Delete Quiz
// @Description Protected Route used to delete a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param quiz_id path string true "Quiz ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/delete/{quiz_id} [DELETE]
func DeleteQuiz(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		err := services.DeleteQuiz(c.Request.Context(), collection, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Quiz Deleted Successfully"})
	}
}



// @Summary Get Quiz
// @Description Protected Route used to get a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param quiz_id path string true "Quiz ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/get/{quiz_id} [GET]
func GetQuiz(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		quiz, err := services.GetQuiz(c.Request.Context(), collection, quizID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quiz)
	}
}



// @Summary Get Quizes By Admin
// @Description Protected Route used to get all quizes by admin
// @Produce json
// @Tags Quizes
// @Accept json
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/admin [GET]
func GetQuizesByAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		quizes, err := services.GetQuizesByAdmin(c.Request.Context(), collection)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizes)
	}
}

// @Summary Get Quizes By Module ID
// @Description Protected Route used to get all quizes by module ID
// @Produce json
// @Tags Quizes
// @Accept json
// @Param module_id path string true "Module ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/module/{module_id} [GET]
func GetQuizesByModuleId(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("Module ID is Required")})
			return
		}
		quizes, err := services.GetQuizesByModuleId(c.Request.Context(), collection, moduleID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizes)
	}
}