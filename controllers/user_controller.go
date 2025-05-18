package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "login-system/config"
    "login-system/middleware"
    "login-system/models"
)

func Register(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 检查用户名是否已存在
    var existingUser models.User
    if result := config.DB.Where("username = ?", input.Username).First(&existingUser); result.RowsAffected > 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
        return
    }

    // 检查邮箱是否已存在
    if result := config.DB.Where("email = ?", input.Email).First(&existingUser); result.RowsAffected > 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
        return
    }

    // 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
        return
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword),
    }

    if result := config.DB.Create(&user); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if result := config.DB.Where("username = ?", input.Username).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 生成JWT
    token, err := middleware.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token生成失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "登录成功",
        "token":   token,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func GetUserProfile(c *gin.Context) {
    userId, _ := c.Get("user_id")
    
    var user models.User
    if result := config.DB.First(&user, userId); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}