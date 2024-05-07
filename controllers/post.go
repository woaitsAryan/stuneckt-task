package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/woaitsAryan/stuneckt-task/cache"
	"github.com/woaitsAryan/stuneckt-task/config"
	"github.com/woaitsAryan/stuneckt-task/helpers"
	"github.com/woaitsAryan/stuneckt-task/initializers"
	"github.com/woaitsAryan/stuneckt-task/models"
	"github.com/woaitsAryan/stuneckt-task/routines"
	"github.com/woaitsAryan/stuneckt-task/schemas"
	API "github.com/woaitsAryan/stuneckt-task/utils"
	"gorm.io/gorm"
)

func GetPost(c *fiber.Ctx) error {
	postID := c.Params("postID")

	postInCache, err := cache.GetPost(postID)

	if err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"post":    postInCache,
		})
	}

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ?", parsedPostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	cache.SetPost(postID, &post)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"post":    post,
	})
}

func GetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("userID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetMyPosts(c *fiber.Ctx) error {
	user := c.Locals("loggedinUser").(*models.User)

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.Preload("User").Where("user_id = ?", user.ID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetMyLikedPosts(c *fiber.Ctx) error {
	user := c.Locals("loggedinUser").(*models.User)

	var postLikes []models.Like
	if err := initializers.DB.Where("user_id = ? AND post_id IS NOT NULL", user.ID).Find(&postLikes).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	var postIDs []string
	for _, post := range postLikes {
		postIDs = append(postIDs, post.PostID.String())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   postIDs,
	})
}

func AddPost(c *fiber.Ctx) error {
	var reqBody schemas.PostCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Req Body"}
	}

	user := c.Locals("loggedinUser").(*models.User)

	newPost := models.Post{
		UserID:  user.ID,
		Content: reqBody.Content,
	}

	result := initializers.DB.Create(&newPost)
	if result.Error != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Added",
		"post":    newPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	user := c.Locals("loggedinUser").(*models.User)

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? and user_id=?", parsedPostID, user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.PostUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Request Body."}
	}

	if reqBody.Content != "" {
		post.Content = reqBody.Content
	}

	post.Edited = true

	if err := initializers.DB.Save(&post).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	cache.RemovePost(postID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	user := c.Locals("loggedinUser").(*models.User)

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? AND user_id=?", parsedPostID, user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}

func LikePost(c *fiber.Ctx) error {
	user := c.Locals("loggedinUser").(*models.User)

	postID := c.Params("postID")
	parsedPostID, err := uuid.Parse(postID)

	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var like models.Like
	err = initializers.DB.Where("user_id=? AND post_id=?", user.ID, parsedPostID).First(&like).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			likeModel := models.Like{
				PostID: &parsedPostID,
				UserID: user.ID,
			}

			result := initializers.DB.Create(&likeModel)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
			go routines.IncrementPostLikes(parsedPostID, user.ID)

		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		result := initializers.DB.Delete(&like)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}
		go routines.DecrementPostLikes(parsedPostID)

	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Liked/Unliked.",
	})
}
