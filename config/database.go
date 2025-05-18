package config

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "login-system/models"
)

var DB *gorm.DB

func ConnectDatabase() {
    dsn := "root:123456789@tcp(127.0.0.1:3306)/login_system?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("数据库连接失败!")
    }

    // 自动迁移
    database.AutoMigrate(&models.User{})

    DB = database
    fmt.Println("数据库连接成功!")
}