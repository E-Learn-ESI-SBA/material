package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"madaurus/dev/material/app/interfaces"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/services"
	"madaurus/dev/material/app/utils"
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
		err = services.ReplayToComment(context.Request.Context(), collection, replay, commentId, user.ID)
		if err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		context.JSON(200, gin.H{"message": "Replay Created Successfully"})
	}
}
