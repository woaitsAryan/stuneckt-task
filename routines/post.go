package routines

import (
	"github.com/woaitsAryan/stuneckt-task/helpers"
	"github.com/woaitsAryan/stuneckt-task/initializers"
	"github.com/woaitsAryan/stuneckt-task/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func IncrementPostLikes(postID uuid.UUID, loggedInUserID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post of this ID found-IncrementPostLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post-IncrementPostLikes", err, "go_routine")
		}
	} else {
		post.NoLikes++

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-IncrementPostLikes", result.Error, "go_routine")
		}
	}
}

func DecrementPostLikes(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Post of this ID found-DecrementPostLikes.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Post-DecrementPostLikes", err, "go_routine")
		}
	} else {
		post.NoLikes--

		result := initializers.DB.Save(&post)
		if result.Error != nil {
			helpers.LogDatabaseError("Error while updating Post-DecrementPostLikes", result.Error, "go_routine")
		}
	}
}