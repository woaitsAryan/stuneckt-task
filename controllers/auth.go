package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/woaitsAryan/stuneckt-task/config"
	"github.com/woaitsAryan/stuneckt-task/helpers"
	"github.com/woaitsAryan/stuneckt-task/initializers"
	"github.com/woaitsAryan/stuneckt-task/models"
	"github.com/woaitsAryan/stuneckt-task/schemas"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {
	var reqBody schemas.UserCreateSchema

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Validation Failed"}
	}	
	
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 12)
	if err != nil {
		go helpers.LogServerError("Error while hashing Password.", err, c.Path())
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.SERVER_ERROR, Err: err}
	}

	newUser := models.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: string(hash),
		Username: reqBody.Username,
		Bio:      reqBody.Bio,
		Links:    reqBody.Links,
	}

	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	c.Set("loggedInUserID", newUser.ID.String())

	access_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": newUser.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.ACCESS_TOKEN_TTL).Unix(),
	})

	access_token, err := access_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		go helpers.LogServerError("Error while decrypting JWT Token.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Account Created",
		"token":   access_token,
		"user":    newUser,
	})

}

func LogIn(c *fiber.Ctx) error {
	var reqBody schemas.UserLoginSchema

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Validation Failed"}
	}

	var user models.User
	if err := initializers.DB.Session(&gorm.Session{SkipHooks: true}).First(&user, "username = ?", reqBody.Username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No account with these credentials found."}
		} else {
			return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No account with these credentials found."}
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	
	access_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.ACCESS_TOKEN_TTL).Unix(),
	})

	access_token, err := access_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		go helpers.LogServerError("Error while decrypting JWT Token.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Signed in!",
		"token":   access_token,
		"user":    user,
	})

}