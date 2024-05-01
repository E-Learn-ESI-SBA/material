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

func CreateComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var comment models.Comments
		err := context.BindJSON(&comment)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.CreateComment(context.Request.Context(), collection, comment)

		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Comment Created Successfully"})

	}

}

func GetCourseComments(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var pagination interfaces.PaginationQuery
		courseId, errP := context.Params.Get("courseId")
		err := context.BindJSON(&pagination)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if errP != true {
			context.JSON(400, gin.H{"error": "CourseId is required"})
			return
		}
		comments, err := services.GetCourseCommentsByCourseId(context.Request.Context(), collection, courseId, pagination)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"comments": comments})
	}
}
func EditComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var comment models.Comments
		err := context.BindJSON(&comment)
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"error": "CommentId is required"})
			return
		}
		user := context.MustGet("user").(utils.UserDetails)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.EditComment(context.Request.Context(), collection, commentId, user.ID, comment)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Comment Updated Successfully"})
	}
}

func EditReplay(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		var comment models.Reply
		err := context.BindJSON(&comment)
		replayId, errP := context.Params.Get("replayId")
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"error": "ReplayId is required"})
			return
		}
		user := context.MustGet("user").(utils.UserDetails)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = services.EditReply(context.Request.Context(), collection, commentId, replayId, user.ID, comment)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Replay Updated Successfully"})
	}
}

func DeleteComment(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"error": "CommentId is required"})
			return
		}
		user := context.MustGet("user").(utils.UserDetails)
		err := services.DeleteCommentByUser(context.Request.Context(), collection, commentId, user.ID)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Comment Deleted Successfully"})
	}
}
func DeleteReplay(collection *mongo.Collection) gin.HandlerFunc {
	return func(context *gin.Context) {
		replayId, errP := context.Params.Get("replayId")
		commentId, errP := context.Params.Get("commentId")
		if errP != true {
			context.JSON(400, gin.H{"error": "ReplayId is required"})
		}
		user := context.MustGet("user").(utils.UserDetails)
		err := services.DeleteReplyByUser(context.Request.Context(), collection, commentId, replayId, user.ID)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
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
			context.JSON(http.StatusBadRequest, gin.H{"error": shared.INVALID_BODY})
			return
		}
		commentId, errP := context.Params.Get("commentId")
		commendObjectId, err := primitive.ObjectIDFromHex(commentId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": shared.INVALID_ID})
			return
		}
		if errP != true {
			context.JSON(400, gin.H{"error": "CommentId is required"})
			return
		}
		value, errC := context.Get("user")
		user := value.(utils.UserDetails)
		if errC != true {
			context.JSON(http.StatusInternalServerError, gin.H{"error": shared.USER_NOT_INJECTED})
			return
		}
		replay.UserId = user.ID
		replay.User = utils.LightUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Avatar:   user.Avatar,
		}

		err = services.ReplayToComment(context.Request.Context(), collection, replay, commendObjectId)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Replay Created Successfully"})
	}
}
