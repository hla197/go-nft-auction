package handlers

import (
	"chain/internal/infra/db"
	"chain/internal/models"
	"chain/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

type RegisterRequest struct {
	*utils.FieldValidate
	Username       string `json:"username" binding:"required,min=3,max=20" label:"用户名"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6"`
}

type LoginRequest struct {
	*utils.FieldValidate
	Username string `json:"username" binding:"required,min=3,max=20" label:"用户名"`
	Password string `json:"password" binding:"required,min=6" label:"密码"`
}

type AuthResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, 500, msg)
		return
	}

	var user models.User
	if err := db.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, 500, "user not found")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, 500, "password incorrect")
		return
	}

	token, err := utils.GenerateToken(uint64(user.ID), user.Username)

	if err != nil {
		log.Panicf("generate token err: %v", err)
		utils.Fail(c, 500, "generate token failed")
		return
	}

	utils.Success(c, &AuthResponse{
		Username: user.Username,
		Token:    token,
	}, "")
	return
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, 500, msg)
		return
	}

	if req.Password != req.Password {
		utils.Fail(c, 500, "two password not match")
		return
	}

	var existUser models.User
	db.GetDB().Where("username = ?", req.Username).First(&existUser)
	if existUser.ID != 0 {
		utils.Fail(c, 500, "username is exist")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash err: %v", err)
		utils.Error(c, "bcrypt password err")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		log.Printf("create user err: %v", err)
		utils.Error(c, "create user fail")
		return
	}

	token, err := utils.GenerateToken(uint64(user.ID), user.Username)

	if err != nil {
		log.Panicf("generate token err: %v", err)
		utils.Fail(c, 500, "generate token failed")
		return
	}
	utils.Success(c, &AuthResponse{
		Username: user.Username,
		Token:    token,
	}, "register success")
	return
}
