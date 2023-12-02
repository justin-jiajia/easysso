package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/justin-jiajia/easysso/api/handler"
	"github.com/justin-jiajia/easysso/api/middleware"
)

func InitApi(r *gin.Engine) {
	api := r.Group("/api/")
	{
		user := api.Group("/user/")
		{
			user.POST("/sign_up/", handler.SignUpHandler) //注册
			user.POST("/sign_in/", handler.SignInHandler) //登录
			setting := user.Group("/settings/")
			setting.Use(middleware.Auth())
			{
				setting.POST("/getcode/", handler.GetCodeHandler)            //获取code
				setting.POST("/avatar/", handler.AvatarUploadHandler)        //上传头像
				setting.POST("/change_passwd/", handler.ChangePasswdHandler) //修改密码
				setting.POST("/remove/", handler.RemoveUserHandler)          //删除账户
			}
		}
		oath2 := api.Group("/oath2/")
		{
			oath2.POST("/getcallback/", handler.GetCallbackHandler) //验证服务器&获取callback
			oath2.POST("/gettoken/", handler.GetTokenHandler)       //验证code
			oath2.POST("/information/", handler.InformationHandler) //在验证后获取详情
		}
		api.GET("/avatar/:id/", handler.AvatarHandler) //查看头像
	}
	api.Use(cors.Default())
}
