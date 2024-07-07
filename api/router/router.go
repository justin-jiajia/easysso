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
			user.POST("/sign_up/", handler.SignUpHandler)                   //注册
			user.POST("/sign_in/", handler.SignInHandler)                   //登录
			user.POST("/startwlogin/", handler.WebauthnLoginStartHandler)   //登录
			user.POST("/finishwlogin/", handler.WebauthnLoginFinishHandler) //登录完成
			setting := user.Group("/settings/")
			setting.Use(middleware.Auth())
			{
				setting.POST("/getcode/", handler.GetCodeHandler)            //获取code
				setting.POST("/avatar/", handler.AvatarUploadHandler)        //上传头像
				setting.POST("/change_passwd/", handler.ChangePasswdHandler) //修改密码
				setting.POST("/remove/", handler.RemoveUserHandler)          //删除账户
				setting.POST("/remove_token/", handler.RemoveTokenHandler)   //删除token
				webauthn := setting.Group("/webauthn/")
				{
					webauthn.GET("/startregistration/", handler.WebauthnResStartHandler)    //注册
					webauthn.POST("/finishregistration/", handler.WebauthnResFinishHandler) //注册完成
					webauthn.GET("/list/", handler.WebauthnList)                            //列举验证器
					webauthn.GET("/log/", handler.GetActions)                               //列举行为日志
					webauthn.POST("/delete/", handler.WebauthnDelete)                       //删除验证器
					webauthn.POST("/edit/", handler.WebauthnEdit)                           //重命名验证器

				}
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
