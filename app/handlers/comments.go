package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"net/http"
)

// @Summary Create Comment
// @Description Protected Route used to create a comment
// @Produce json
// @Accept json
// @Tags Comments
// @Param comment body models.Comments true "Comment Object"
// @Param courseId query string true "Course ID"
// @Success 201 {object} interfaces.APIResponse
// @Failure 400 {object} interfaces.APIResponse
// @Failure 500 {object} interfaces.APIResponse
// @Router /comments [POST]
// @Security Bearer
func CreateComment(collection *mongo.Collection, userCollection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		courseId := context.Query("courseId")
		if courseId == "" {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		courseObjectId, err := primitive.ObjectIDFromHex(courseId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}

		var comment models.Comments
		err = context.ShouldBindJSON(&comment)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		comment.CourseId = courseObjectId
		value, errU := context.Get("user")
		if errU != true {
			context.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		user := value.(*utils.UserDetails)
		comment.UserId = user.ID
		err = services.CreateComment(context.Request.Context(), collection, comment, userCollection, *user)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusCreated, gin.H{"message": shared.COMMENT_CREATED})
	}

}

func GetCourseComments(collection *mongo.Collection, userCollection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var pagination interfaces.PaginationQuery
		courseId, errP := context.Params.Get("courseId")
		err := context.ShouldBindJSON(&pagination)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		if errP != true {
			context.JSON(http.StatusBadRequest, gin.H{"message": "CourseId is required"})
			return
		}
		comments, err := services.GetCourseCommentsByCourseId(context.Request.Context(), collection, courseId, pagination, userCollection)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		context.JSON(200, gin.H{"data": comments})
	}
}
func EditComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var comment models.Comments
		err := context.BindJSON(&comment)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": shared.INVALID_BODY})
			return
		}
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.REQUIRED_ID})
			return
		}
		user, err := utils.GetUserPayload(context)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = services.EditComment(context.Request.Context(), collection, commentId, user.ID, comment)
		if err != nil {
			context.JSON(400, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": shared.COMMENT_UPDATED})
	}
}

func EditReplay(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var comment models.Reply
		replayId, errP := context.Params.Get("replayId")
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(http.StatusBadRequest, gin.H{"message": "ReplayId is required"})
			return
		}
		user, err := utils.GetUserPayload(context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = context.BindJSON(&comment)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = services.EditReply(context.Request.Context(), collection, commentId, replayId, user.ID, comment)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": shared.COMMENT_UPDATED})
	}
}

func DeleteComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"message": "CommentId is required"})
			return
		}
		user, err := utils.GetUserPayload(context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = services.DeleteCommentByUser(context.Request.Context(), collection, commentId, user.ID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": "Comment Deleted Successfully"})
	}
}
func DeleteReplay(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		replayId, errP := context.Params.Get("replayId")
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"message": "ReplayId is required"})
		}
		user, err := utils.GetUserPayload(context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = services.DeleteReplyByUser(context.Request.Context(), collection, commentId, replayId, user.ID)
		if err != nil {
			context.JSON(400, gin.H{"message": err.Error()})
			return

		}
		context.JSON(200, gin.H{"message": "Replay Deleted Successfully"})
	}
}

func ReplayToComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var replay models.Reply
		err := context.BindJSON(&replay)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_BODY})
			return
		}
		commentId, errP := context.Params.Get("commentId")
		commendObjectId, err := primitive.ObjectIDFromHex(commentId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": shared.INVALID_ID})
			return
		}
		if errP != true {
			context.JSON(400, gin.H{"message": "CommentId is required"})
			return
		}
		value, errC := context.Get("user")
		user := value.(*utils.UserDetails)
		if errC != true {
			context.JSON(http.StatusInternalServerError, gin.H{"message": shared.USER_NOT_INJECTED})
			return
		}
		replay.UserId = user.ID

		err = services.ReplayToComment(context.Request.Context(), collection, replay, commendObjectId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusCreated, gin.H{"message": "Replay Created Successfully"})
	}
}
