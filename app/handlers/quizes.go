package handlers

import (
	"errors"
	"log"
	"madaurus/dev/material/app/kafka"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Router /quizes [POST]
func CreateQuiz(collection *mongo.Collection, moduleCollection *mongo.Collection, instance *kafka.KafkaInstance) gin.HandlerFunc {
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
		go instance.ResourceCreatingProducer(user, "Quiz", quiz.Title, kafka.USER_NOTIFICATION_TYPE)
		c.JSON(200, gin.H{"message": shared.QUIZ_CREATED})
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
		var quizUpdates models.QuizUpdate
		user := c.MustGet("user").(*utils.UserDetails)
		err := c.BindJSON(&quizUpdates)
		log.Println(quizUpdates)
		log.Println("handler update quiz")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		err = services.UpdateQuiz(c.Request.Context(), collection, quizUpdates, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": shared.QUIZ_UPDATED})
	}

}



// @Summary Delete Quiz
// @Description Protected Route used to delete a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param id path string true "Quiz ID"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/delete/{id} [DELETE]
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
		c.JSON(200, gin.H{"message": shared.QUIZ_DELETED})
	}
}



// @Summary Get Quiz
// @Description Protected Route used to get a quiz
// @Produce json
// @Tags Quizes
// @Accept json
// @Param id path string true "Quiz ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/{id} [GET]
func GetQuiz(collection *mongo.Collection, submissionCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		quizID, errP := c.Params.Get("id")
		user := c.MustGet("user").(*utils.UserDetails)
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		quiz, err, passed := services.GetQuiz(c.Request.Context(), collection, submissionCollection, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"quiz": quiz, "passed": passed})
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
// @Param id path string true "Module ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/module/{id} [GET]
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

// @Summary Get Many Quizes By Teacher ID
// @Description Protected Route used to get all quizes by teacher ID
// @Produce json
// @Tags Quizes
// @Accept json
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/teacher [GET]
func GetManyQuizesByTeacherId(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		if (user.Role != "teacher") {
			c.JSON(400, gin.H{"error": errors.New("Only Teachers can access this route")})
			return
		}
		quizes, err := services.GetManyQuizesByTeacherId(c.Request.Context(), collection, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizes)
	}
}


// @Summary Get Quizes By Student ID
// @Description Protected Route used to get all quizes by student ID
// @Produce json
// @Tags Quizes
// @Accept json
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/student [GET]
func GetQuizesByStudentId(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizes, err := services.GetQuizesByStudentId(c.Request.Context(), collection, user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizes)
	}
}


// @Summary Submit Quiz Answers
// @Description Protected Route used to submit quiz answers by a student
// @Produce json
// @Tags Quizes
// @Accept json
// @Param id path string true "Quiz ID"
// @Param answer body models.QuizAnswer true "Quiz Answer Object"
// @Success 200 {object} interfaces.APiSuccess
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/submit/{id} [POST]
func SubmitQuizAnswers(collection *mongo.Collection, SubmissionsCollection *mongo.Collection, instance *kafka.KafkaInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		var submission models.Submission
		err := c.BindJSON(&submission)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		submission.QuizId, err = primitive.ObjectIDFromHex(quizID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		submission.StudentId = user.ID
		score, err := services.SubmitQuizAnswers(c.Request.Context(), collection, SubmissionsCollection, submission)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		go instance.EvaluationProducer(user, "Quiz", score)
		c.JSON(200, gin.H{"message": "Quiz Answer Submitted Successfully"})
	}
}



// @Summary Get Quiz Results By teacher
// @Description Protected Route used to get quiz results by teacher
// @Produce json
// @Tags Quizes Submissions
// @Accept json
// @Param id path string true "Quiz ID"
// @Success 200 {object} models.QuizResults
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/{id}/results [GET]
func GetQuizResults(collection *mongo.Collection, SubmissionsCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		quizResults, err := services.GetQuizResults(c.Request.Context(), collection, SubmissionsCollection, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizResults)
	}
}

// @Summary Get Submission Details by submission id (teacher only endpoint)
// @Description Protected Route used to get submission details by submission id
// @Produce json
// @Tags Quizes Submissions
// @Accept json
// @Param id path string true "Submission ID"
// @Success 200 {object} models.Submission
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/:id/teacher/result [GET]
func GetSubmissionDetails(collection *mongo.Collection, SubmissionsCollection *mongo.Collection, moduleCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		submissionID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("Submission ID is Required")})
			return
		}
		if (user.Role != "teacher") {
			c.JSON(400, gin.H{"error": errors.New("Only Teachers can access this route")})
			return
		}
		submission, module_name, quiz, err := services.GetSubmissionDetails(c.Request.Context(), collection, moduleCollection, SubmissionsCollection, submissionID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(200, gin.H{"submission": submission, "module_name": *module_name, "quiz": quiz})
	}
}



// @Summary Get Quiz Result By Student ID
// @Description Protected Route used to get quiz result by student ID
// @Produce json
// @Tags Quizes Submissions
// @Accept json
// @Success 200 {object} models.QuizResults
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/:id/student/result [GET]
func GetQuizResultByStudentId(
	collection *mongo.Collection, 
	SubmissionsCollection *mongo.Collection,
	moduleCollection *mongo.Collection,
	) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		submission, quiz, module_name, err := services.GetQuizResultByStudentId(c.Request.Context(), collection, SubmissionsCollection, moduleCollection, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"submission": submission, "quiz": quiz, "module_name": *module_name})
	}
}

// @Summary Get Quizes Results By Student ID
// @Description Protected Route used to get all quizes results by student ID
// @Produce json
// @Tags Quizes Submissions
// @Accept json
// @Success 200 {object} models.QuizResults
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/results [GET]
func GetQuizesResultsByStudentId(collection *mongo.Collection, SubmissionsCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizResults, err := services.GetQuizesResultsByStudentId(c.Request.Context(), collection, SubmissionsCollection, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quizResults)
	}
}

// @Summary Get Quiz Questions By student
// @Description Protected Route used to get quiz questions by student
// @Produce json
// @Tags Quizes
// @Accept json
// @Param id path string true "Quiz ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} interfaces.APiError
// @Router /quizes/:id/questions [GET]
func GetQuizQuestions(collection *mongo.Collection, submissionCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*utils.UserDetails)
		quizID, errP := c.Params.Get("id")
		if !errP {
			c.JSON(400, gin.H{"error": errors.New("quiz ID is Required")})
			return
		}
		quiz, err := services.GetQuizQuestions(c.Request.Context(), collection, submissionCollection, quizID, user.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"questions": quiz.Questions,
			"duration": quiz.Duration,
			"title": quiz.Title,
		})
	}
}