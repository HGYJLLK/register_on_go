package routes

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "login-system/controllers"
    "login-system/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:5500"},  // 添加HTML测试页面的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},  // 添加OPTIONS方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

    // 公开路由
    public := r.Group("/api")
    {
        public.POST("/register", controllers.Register)
        public.POST("/login", controllers.Login)
    }

    // 需要授权的路由
    protected := r.Group("/api")
    protected.Use(middleware.JWTAuth())
    {
        protected.GET("/user/profile", controllers.GetUserProfile)
    }

    return r
}