package main

import (
    "login-system/config"
    "login-system/routes"
)

func main() {
    // 连接数据库
    config.ConnectDatabase()

    // 设置路由
    r := routes.SetupRouter()

    // 启动服务器
    r.Run(":8000")
}